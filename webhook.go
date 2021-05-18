package wecomrus

import (
	"time"

	"github.com/juju/ratelimit"
)

type WebHook struct {
	bucket   *ratelimit.Bucket
	endpoint string
}

var webhooks []WebHook

func newWebhook(endpoint string) {
	webhooks = append(webhooks, WebHook{
		bucket:   ratelimit.NewBucketWithQuantum(time.Minute, 20, 20),
		endpoint: endpoint,
	})
}
