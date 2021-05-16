package wecomrus

import (
	"time"
)

// Option defines ...
type Option struct {
	// AppName will be used for logging out
	AppName string
	// MessageFormart defines the format of the message
	MessageFormart string
	// ChatID defines ...
	ChatID string
	// TimeFormat defines ...
	TimeFormat string
	// TimeZone defines ...
	TimeZone *time.Location

	// CorpID defines ...
	CorpID string
	// CorpSecret defines ...
	CorpSecret string
}

var options Option

func mergeOptions(opts ...Option) {
	for _, opt := range opts {
		if opt.AppName != "" {
			options.AppName = opt.AppName
		}
		if opt.MessageFormart != "" {
			options.MessageFormart = opt.MessageFormart
		}
		if opt.ChatID != "" {
			options.ChatID = opt.ChatID
		}
		if opt.CorpSecret != "" {
			options.CorpSecret = opt.CorpSecret
		}
		if opt.CorpID != "" {
			options.CorpID = opt.CorpID
		}
		if opt.TimeFormat != "" {
			options.TimeFormat = opt.TimeFormat
		}
		if opt.TimeZone != nil {
			options.TimeZone = opt.TimeZone
		}
	}
}

func checkOptions() {
	if options.AppName == "" {
		options.AppName = "Undefined"
	}
	if options.TimeZone == nil {
		tz, err := time.LoadLocation("Asia/Chongqing")
		if err != nil {
			panic(err)
		}
		options.TimeZone = tz
	}
	if options.TimeFormat == "" {
		options.TimeFormat = "01-02 15:04:05"
	}
	if options.MessageFormart == "" {
		// TODO:---
		options.MessageFormart =
			`Log Monitor
***********
* AppName: {{app}}
* Time: {{time}}
* Level: {{level}}
***********
* Message: {{message}}
***********
{{content}}`
	}
}
