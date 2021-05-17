package wecomrus

import (
	"log"
	"time"
)

// SafeSwitcher defines a switcher for switching safe flag
type SafeSwitcher bool

const (
	// SafeOn enables the message safe
	SafeOn SafeSwitcher = true
	// SafeOff disenables the message safe
	SafeOff SafeSwitcher = false
)

// ToInt for translating to int
func (s SafeSwitcher) ToInt() int {
	if s {
		return 1
	}

	return 0
}

// MessageType defines ...
type MessageType string

const (
	// TextMessage for sending text message
	TextMessage MessageType = "text"
)

type MessageFormater interface {
	ToString() string
}

type normalTextFormat struct {
	text string
}

func (n normalTextFormat) ToString() string {
	return n.text
}

// Option defines ...
type Option struct {
	// AppName will be used for logging out
	AppName string
	// MessageFormart defines the format of the message
	// It is an interface that you can customize the output as long as you implement it
	MessageFormat MessageFormater
	// GroupChatID defines ...
	GroupChatID string
	// TimeFormat defines ...
	TimeFormat string
	// TimeZone defines ...
	TimeZone *time.Location
	// Safe defines ...
	Safe SafeSwitcher
	// MsgType defines what type of log you want to display in WeCom
	MsgType MessageType

	// CorpID defines ...
	CorpID string
	// CorpSecret defines ...
	CorpSecret string
}

var options Option

func mergeOptions(opts ...Option) {
	// Set defaults
	options.MsgType = TextMessage

	for _, opt := range opts {
		if opt.AppName != "" {
			options.AppName = opt.AppName
		}
		if opt.MessageFormat != nil {
			options.MessageFormat = opt.MessageFormat
		}
		if opt.GroupChatID != "" {
			options.GroupChatID = opt.GroupChatID
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
		if opt.Safe == SafeOn {
			options.Safe = SafeOn
		}
		if opt.MsgType != TextMessage {
			options.MsgType = opt.MsgType
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
	if options.MessageFormat == nil {
		var normalMessage normalTextFormat
		normalMessage.text = `WeComRus Monitor
***********
* AppName: {{app}}
* Time: {{time}}
* Level: {{level}}
***********
* Message: {{message}}
***********
{{content}}`
		options.MessageFormat = normalMessage
	}
	if options.CorpID == "" {
		log.Println("Warning: You have not defined the CorpID")
	}
	if options.CorpSecret == "" {
		log.Println("Warning: You have not defined the CorpSecret")
	}
	if options.GroupChatID == "" {
		log.Println("Warning: You have not provided any group chat IDs in the GroupChatIDs")
	}
}
