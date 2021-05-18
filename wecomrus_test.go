package wecomrus

import (
	"testing"

	"github.com/henry0475/wecomrus/options"
	"github.com/sirupsen/logrus"
)

func TestWeComRus(t *testing.T) {
	hook, err := NewWeComHook(
		&options.Option{
			CorpID:      "",
			CorpSecret:  "",
			GroupChatID: "",
			AppName:     "Test APP",
		},
		&options.Option{
			Safe: options.SafeOff,
		},
	)
	if err != nil {
		t.Log(err.Error())
	}

	log := logrus.New()

	if err == nil {
		log.Hooks.Add(hook)
	}
	log.WithFields(logrus.Fields{
		"t1": "a1",
	}).Warn("asass")
	log.Error("1234567")

	t.Error("Done")
}
