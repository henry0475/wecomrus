package options

import (
	"log"
	"time"
)

// Bool returns a pointer to the string value passed in.
func Bool(b bool) *bool {
	return &b
}

var True = Bool(true)
var False = Bool(false)

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
	// MarkdownMessage for sending markdown message
	// MarkdownMessage MessageType = "markdown"
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

// option defines ...
type Option struct {
	// AppName will be used for logging out
	AppName string
	// MessageFormart defines the format of the message
	// It is an interface that you can customize the output as long as you implement it
	MessageFormat MessageFormater
	// GroupChatID defines ...
	GroupChatID string
	// Webhooks defines ...
	Webhooks []string
	// EnableStats determine whether statistics has enabled
	EnableStats *bool
	// TimeFormat defines ...
	TimeFormat string
	// TimeZone defines ...
	TimeZone *time.Location
	// Safe defines ...
	Safe SafeSwitcher
	// MsgType defines what type of log you want to display in WeCom
	MsgType MessageType
	// DropIfFull defines ...
	DropIfFull *bool

	// CorpID defines ...
	CorpID string
	// CorpSecret defines ...
	CorpSecret string
}

var options Option

func MergeOptions(opts ...*Option) {
	// Set defaults
	options.MsgType = TextMessage
	options.EnableStats = True
	options.DropIfFull = True

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AppName != "" {
			options.AppName = opt.AppName
		}
		if opt.MessageFormat != nil {
			options.MessageFormat = opt.MessageFormat
		}
		if opt.GroupChatID != "" {
			options.GroupChatID = opt.GroupChatID
		}
		if len(opt.Webhooks) != 0 {
			options.Webhooks = opt.Webhooks
		}
		if opt.EnableStats != nil {
			options.EnableStats = opt.EnableStats
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
		if opt.MsgType != "" {
			options.MsgType = opt.MsgType
		}
	}

	checkOptions()
}

func checkOptions() {
	if options.AppName == "" {
		options.AppName = "Undefined"
	}
	if options.TimeZone == nil {
		tz, err := time.LoadLocation("Asia/Chongqing")
		if err != nil {
			log.Printf("Warning: we don't find the NZ data, msg: %v\n", err)
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
	if options.GroupChatID == "" && len(options.Webhooks) == 0 {
		log.Println("Warning: You have not provided any group chat IDs or webhooks")
	}
}

// GetOptions return ...
func GetOptions() Option {
	return options
}
