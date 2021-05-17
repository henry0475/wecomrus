package wecomrus

import (
	"testing"
	"time"
)

type markdownMessage struct {
	content string
}

func (m markdownMessage) ToString() string {
	return m.content
}

func Test_mergeOptions(t *testing.T) {
	type args struct {
		opts []Option
	}
	loc, _ := time.LoadLocation("America/New_York")
	var md markdownMessage
	md.content = "##aaa"

	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
		{
			name: "Test1", args: args{
				opts: []Option{
					{
						AppName:     "app1",
						GroupChatID: "c",
					},
					{
						CorpID:     "a",
						CorpSecret: "b",
					},
					{
						TimeFormat: "2006-01-02 15:04:05",
						TimeZone:   loc,
					},
					{
						MessageFormat: md,
					},
					{
						Safe: SafeOn,
					},
					{
						MsgType: TextMessage,
					},
				},
			}, want: Option{
				AppName:       "app1",
				CorpID:        "a",
				CorpSecret:    "b",
				GroupChatID:   "c",
				TimeFormat:    "2006-01-02 15:04:05",
				TimeZone:      loc,
				MessageFormat: md,
				Safe:          SafeOn,
				MsgType:       TextMessage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeOptions(tt.args.opts...)
			if options.AppName != tt.want.AppName {
				t.Errorf("AppName is not assigned, got: %v; want: %v", options.AppName, tt.want.AppName)
			}
			if options.CorpID != tt.want.CorpID {
				t.Errorf("CorpID is not assigned, got: %v; want: %v", options.CorpID, tt.want.CorpID)
			}
			if options.CorpSecret != tt.want.CorpSecret {
				t.Errorf("CorpSecret is not assigned, got: %v; want: %v", options.CorpSecret, tt.want.CorpSecret)
			}
			if options.GroupChatID != tt.want.GroupChatID {
				t.Errorf("GroupChatID is not assigned, got: %v; want: %v", options.GroupChatID, tt.want.GroupChatID)
			}
			if options.TimeFormat != tt.want.TimeFormat {
				t.Errorf("TimeFormat is not assigned, got: %v; want: %v", options.TimeFormat, tt.want.TimeFormat)
			}
			if options.TimeZone != tt.want.TimeZone {
				t.Errorf("TimeZone has error, got: %v; want: %v", options.TimeZone, tt.want.TimeZone)
			}
			if options.MessageFormat != tt.want.MessageFormat {
				t.Errorf("MessageFormat has error, got: %v; want: %v", options.MessageFormat, tt.want.MessageFormat)
			}
			if options.Safe != tt.want.Safe {
				t.Errorf("Safe has error, got: %v; want: %v", options.Safe, tt.want.Safe)
			}
		})
	}
}
