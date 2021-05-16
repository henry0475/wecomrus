package wecomrus

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWeComRus(t *testing.T) {
	hook, err := NewWeComHook(Option{
		CorpID:     "",
		CorpSecret: "",
		ChatID:     "",
		AppName:    "Test APP",
	})
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
