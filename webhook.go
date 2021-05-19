package wecomrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/henry0475/wecomrus/options"
	"github.com/henry0475/wecomrus/storage"
	"github.com/henry0475/wecomrus/utils/hash"
	"github.com/juju/ratelimit"
)

type hook struct {
	bucket   *ratelimit.Bucket
	endpoint string
	c        *http.Client
}

type hooks []hook

// Send defines ...
func (h hooks) Send(message string) error {
	for _, hk := range h {
		if hk.bucket.TakeAvailable(1) != 0 {
			// Allowed to send
			storage.Counter.AfterFired(
				hash.GetDestID(hk.getEndPoint()),
				hk.fire(message),
			)
			break
		}
	}
	return nil
}

var webhooks hooks

func loadWebhooks(client *http.Client) {
	for _, wh := range options.GetOptions().Webhooks {
		webhooks = append(webhooks, hook{
			bucket:   ratelimit.NewBucketWithQuantum(time.Minute, 20, 20),
			endpoint: wh,
			c:        client,
		})
	}
}

func (h hook) getEndPoint() string {
	return h.endpoint
}

func (h hook) fire(message string) error {
	var request struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content             string   `json:"content"`
			MentionedList       []string `json:"mentioned_list"`
			MentionedMobileList []string `json:"mentioned_mobile_list"`
		} `json:"text"`
	}
	request.MsgType = string(options.GetOptions().MsgType)
	request.Text.Content = message
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.getEndPoint(), bytes.NewReader(requestJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := h.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response struct {
		Errcode int64  `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	json.Unmarshal(body, &response)

	if response.Errcode == 0 && response.Errmsg == "ok" {
		return nil
	}

	return fmt.Errorf("error code %d with message: %s", response.Errcode, response.Errmsg)
}
