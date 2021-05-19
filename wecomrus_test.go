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
			Webhooks: []string{"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=fe751bb1-e0b8-45cf-8d2d-20156a25b70f"},
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
	time.Sleep(time.Second * time.Duration(2))

	log.Println(storage.Counter.GetStat(hash.GetDestID("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=fe751bb1-e0b8-45cf-8d2d-20156a25b70f")))

	t.Error("Done")
}
