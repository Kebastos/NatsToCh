package cache

import (
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	maxSize     int
	maxWait     time.Duration
	maxByteSize int
	items       []interface{}
	cfg         *config.BufferConfig
	c           chan []interface{}
}

func New(cfg *config.BufferConfig, c chan []interface{}) *Cache {
	items := make([]interface{}, 0)

	cache := Cache{
		items:       items,
		maxSize:     cfg.MaxSize,
		maxWait:     cfg.MaxWait,
		maxByteSize: cfg.MaxByteSize,
		cfg:         cfg,
		c:           c,
	}

	cache.startCleaner()

	return &cache
}

func (c *Cache) Set(value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.items = append(c.items, value)
}

func (c *Cache) Count() int {
	return len(c.items)
}

func (c *Cache) startCleaner() {
	go c.clean(func() bool { return len(c.items) > c.maxSize }, "size")
	go c.clean(func() bool { return cap(c.items) > c.maxByteSize }, "byte size")
	go c.clean(func() bool { <-time.After(c.maxWait); return true }, "time")
}

func (c *Cache) clean(condition func() bool, cleanType string) {
	for {
		if condition() {
			log.Debugf("cleanup by %s complete", cleanType)
			c.c <- c.items
			c.items = make([]interface{}, 0)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
