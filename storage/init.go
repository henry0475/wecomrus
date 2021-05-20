package storage

import (
	"time"
)

func init() {
	data.data = make(map[string]interface{})
	data.setObject(startupTime, time.Now().UnixNano())

	Counter.successMap = make(map[string]requestCounter)
	Counter.failMap = make(map[string]requestCounter)
}
