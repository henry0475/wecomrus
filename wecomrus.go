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

	"github.com/henry0475/wecomrus/options"
	"github.com/henry0475/wecomrus/storage"
	"github.com/henry0475/wecomrus/tokens"
	"github.com/henry0475/wecomrus/utils/hash"
	"github.com/sirupsen/logrus"
)

// WeComHook is ...
type WeComHook struct {
	t tokens.Token
	c http.Client
}

func NewWeComHook(opts ...*options.Option) (*WeComHook, error) {
	options.MergeOptions(opts...)

	return &WeComHook{
		t: tokens.NewAccessToken(options.GetOptions().CorpID, options.GetOptions().CorpSecret),
		c: http.Client{},
	}, nil
}

func getMessage(entry *logrus.Entry) string {
	message := options.GetOptions().MessageFormat.ToString()
	message = strings.ReplaceAll(message, "{{app}}", options.GetOptions().AppName)
	message = strings.ReplaceAll(message, "{{time}}", entry.Time.In(options.GetOptions().TimeZone).Format(options.GetOptions().TimeFormat))
	message = strings.ReplaceAll(message, "{{level}}", entry.Level.String())
	message = strings.ReplaceAll(message, "{{message}}", entry.Message)
	fields, _ := json.MarshalIndent(entry.Data, "", "\t")
	message = strings.ReplaceAll(message, "{{content}}", string(fields))
	return message
}

func (hook *WeComHook) sendToGroupChat(ctx context.Context, message string) error {
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

func (hook *WeComHook) sendByWebhooks(ctx context.Context, message string) error {

	// TODO: xxx
	var _ = func(ctx context.Context, message string) {

	}

	for _, webhook := range options.GetOptions().Webhooks {

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

		req, err := http.NewRequestWithContext(ctx, "POST", webhook, bytes.NewReader(requestJSON))
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

	return nil
}

// Fire is called when a log event is fired.
func (hook *WeComHook) Fire(entry *logrus.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	go func() {
		if options.GetOptions().GroupChatID == "" {
			return
		}
		if err := hook.sendToGroupChat(ctx, getMessage(entry)); err != nil {
			if options.GetOptions().EnableStats == options.Bool(true) {
				storage.Counter.LogFailedToSend(hash.GetDestID(options.GetOptions().GroupChatID))
			}
		} else {
			if options.GetOptions().EnableStats == options.Bool(true) {
				storage.Counter.LogSentTo(hash.GetDestID(options.GetOptions().GroupChatID))
			}
		}
	}()

	return nil
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
