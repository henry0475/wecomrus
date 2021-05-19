package wecomrus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/henry0475/wecomrus/options"
	"github.com/henry0475/wecomrus/tokens"
)

type Sender interface {
	Send(ctx context.Context, message string) error
}

type GroupChat struct {
	t        tokens.Token
	c        *http.Client
	endpoint string
}

var groupChat GroupChat

func loadGroupChat(client *http.Client) {
	groupChat.t = tokens.NewAccessToken(options.GetOptions().CorpID, options.GetOptions().CorpSecret)
	groupChat.c = client
	groupChat.endpoint = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send"
}

func (g *GroupChat) Send(ctx context.Context, message string) error {
	if options.GetOptions().GroupChatID == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(2)*time.Second)
	defer cancel()

	var request struct {
		ChatID  string `json:"chatid"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		IsSafe int `json:"safe"`
	}
	request.ChatID = options.GetOptions().GroupChatID
	request.MsgType = string(options.GetOptions().MsgType)
	request.Text.Content = message
	request.IsSafe = options.GetOptions().Safe.ToInt()
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.endpoint+"?access_token="+g.t.ToString(), bytes.NewReader(requestJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := g.c.Do(req)
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
