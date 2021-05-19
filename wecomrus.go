package wecomrus

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/henry0475/wecomrus/options"
	"github.com/sirupsen/logrus"
)

// WeComHook is ...
type WeComHook struct {
	senders []Sender
}

func NewWeComHook(opts ...*options.Option) (*WeComHook, error) {
	options.MergeOptions(opts...)
	client := &http.Client{}
	// Try to load
	loadWebhooks(client)
	loadGroupChat(client)

	wch := &WeComHook{}

	if len(webhooks) != 0 {
		wch.senders = append(wch.senders, webhooks)
	}
	if options.GetOptions().GroupChatID != "" {
		wch.senders = append(wch.senders, &groupChat)
	}

	return wch, nil
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

// Fire is called when a log event is fired.
func (hook *WeComHook) Fire(entry *logrus.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	message := getMessage(entry)
	for _, sender := range hook.senders {
		go sender.Send(ctx, message)
	}
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
