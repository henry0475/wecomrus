package wecomrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/henry0475/wecomrus/options"
	"github.com/juju/ratelimit"
)

type hooker interface {
	getEngPoint() string
}

type hook struct {
	bucket   *ratelimit.Bucket
	endpoint string
	c        *http.Client
}

type hooks []hook

func (h hooks) send() {
	for _, hk := range h {
		hk.bucket.TakeAvailable(1)
	}
}

var webhooks hooks

func newWebhook(client *http.Client, endpoint string) {
	webhooks = append(webhooks, hook{
		bucket:   ratelimit.NewBucketWithQuantum(time.Minute, 20, 20),
		endpoint: endpoint,
		c:        client,
	})
}

func (hook webHook) getEngPoint() string {
	return hook.endpoint
}

func (hook webHook) fire(message string) error {
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

	req, err := http.NewRequest("POST", hook.getEngPoint(), bytes.NewReader(requestJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := hook.c.Do(req)
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
