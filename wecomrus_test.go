package wecomrus

import (
	"testing"
	"time"

	"github.com/henry0475/wecomrus/options"
	"github.com/henry0475/wecomrus/storage"
	"github.com/henry0475/wecomrus/utils/hash"
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
			Webhooks: []string{""},
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

	for i := 0; i < 30; i++ {
		log.WithFields(logrus.Fields{
			"for": "test limit",
		}).Error(i)
	}

	time.Sleep(time.Second * time.Duration(30))

	log.Println(storage.Counter.GetStat(hash.GetDestID("")))

	t.Error("Done")
}
