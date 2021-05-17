package wecomrus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/henry0475/wecomrus/tokens"
	"github.com/sirupsen/logrus"
)

// WeComHook is ...
type WeComHook struct {
	t tokens.Token
	c http.Client
}

func NewWeComHook(opts ...Option) (*WeComHook, error) {
	mergeOptions(opts...)
	checkOptions()

	return &WeComHook{
		t: tokens.NewAccessToken(options.CorpID, options.CorpSecret),
		c: http.Client{},
	}, nil
}

func getMessage(entry *logrus.Entry) string {
	message := options.MessageFormat.ToString()
	message = strings.ReplaceAll(message, "{{app}}", options.AppName)
	message = strings.ReplaceAll(message, "{{time}}", entry.Time.In(options.TimeZone).Format(options.TimeFormat))
	message = strings.ReplaceAll(message, "{{level}}", entry.Level.String())
	message = strings.ReplaceAll(message, "{{message}}", entry.Message)
	fields, _ := json.MarshalIndent(entry.Data, "", "\t")
	message = strings.ReplaceAll(message, "{{content}}", string(fields))
	return message
}

// Fire is called when a log event is fired.
func (hook *WeComHook) Fire(entry *logrus.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	var request struct {
		ChatID  string `json:"chatid"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		IsSafe int `json:"safe"`
	}
	request.ChatID = options.GroupChatID
	request.MsgType = string(options.MsgType)
	request.Text.Content = getMessage(entry)
	request.IsSafe = options.Safe.ToInt()
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token="+hook.t.ToString(), bytes.NewReader(requestJSON))
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

// Levels returns the available logging levels.
func (hook *WeComHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}
