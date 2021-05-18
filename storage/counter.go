package storage

import (
	"sync"
	"time"

	"github.com/henry0475/wecomrus/options"
)

type requestCounter map[string]int

type CountCollection struct {
	successMap map[string]requestCounter
	failMap    map[string]requestCounter
	mu         sync.RWMutex
}

var Counter CountCollection

// LogSentTo will ...
func (c *CountCollection) LogSentTo(destID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	date := time.Now().In(options.GetOptions().TimeZone).Format("2006-01-02")
	if curr, ok := c.successMap[date][destID]; ok {
		c.successMap[date][destID] = curr + 1
	} else {
		c.successMap[date][destID] = 1
	}
}

// LogFailedToSend will ...
func (c *CountCollection) LogFailedToSend(destID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	date := time.Now().In(options.GetOptions().TimeZone).Format("2006-01-02")
	if curr, ok := c.failMap[date][destID]; ok {
		c.failMap[date][destID] = curr + 1
	} else {
		c.failMap[date][destID] = 1
	}
}

// GetSentCount returns ...
// date => 2006-01-02
func (c *CountCollection) GetSentCount(date string, destID string) int {
	if date == "" {
		date = time.Now().In(options.GetOptions().TimeZone).Format("2006-01-02")
	}
	if curr, ok := c.successMap[date][destID]; ok {
		return curr
	}
	return 0
}

// GetUnsentCount returns ...
// date => 2006-01-02
func (c *CountCollection) GetUnsentCount(date string, destID string) int {
	if date == "" {
		date = time.Now().In(options.GetOptions().TimeZone).Format("2006-01-02")
	}
	if curr, ok := c.failMap[date][destID]; ok {
		return curr
	}
	return 0
}

// Stat defines ...
type Stat struct {
	Sent     int
	Unsent   int
	Duration int64 // second
}

// GetStat returns ...
func (c *CountCollection) GetStat(destID string) Stat {
	var sentCount, unsentCount int
	for _, val := range c.successMap {
		if count, ok := val[destID]; ok {
			sentCount = sentCount + count
		}
	}
	for _, val := range c.failMap {
		if count, ok := val[destID]; ok {
			unsentCount = unsentCount + count
		}
	}

	return Stat{
		Sent:     sentCount,
		Unsent:   unsentCount,
		Duration: GetUpTime(),
	}
}