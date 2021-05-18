package storage

import (
	"fmt"
	"sync"
	"time"
)

type dataStorage struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var data dataStorage

func (d *dataStorage) setObject(name string, val interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[name] = val
}

func (d *dataStorage) getObject(name string) (interface{}, error) {
	if v, ok := d.data[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("no data found with name %s", name)
}

// GetStartupTime returns the startup time in nanoseconds
func GetStartupTime() int64 {
	if obj, err := data.getObject(startupTime); err == nil {
		if t, ok := obj.(int64); ok {
			return t
		}
	}
	return 0
}

// GetUpTime returns the uptime in seconds
func GetUpTime() int64 {
	if obj, err := data.getObject(startupTime); err == nil {
		if t, ok := obj.(int64); ok {
			return time.Now().Unix() - time.Unix(0, t).Unix()
		}
	}
	return 0
}
