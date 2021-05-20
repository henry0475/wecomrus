package wecomrus

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/henry0475/wecomrus/options"
	"github.com/sirupsen/logrus"
)

// WeComHook is ...
type WeComHook struct {
	senders []Sender
}

// NewWeComHook for creating the hook
func NewWeComHook(opts ...*options.Option) (*WeComHook, error) {
	options.MergeOptions(opts...)
	client := &http.Client{}
	loadWebhooks(client)
	if options.GetOptions().CorpID != "" && options.GetOptions().CorpSecret != "" {
		loadGroupChat(client)
	}

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
	message := getMessage(entry)
	for _, sender := range hook.senders {
		go func(sender Sender) {
			if err := sender.Send(message); err != nil {
				log.Println(err)
			}
		}(sender)
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
