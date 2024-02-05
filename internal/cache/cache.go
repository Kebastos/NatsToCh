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
	go c.cleanBySize()
	go c.cleanByteSize()
	go c.cleanByTime()
}

func (c *Cache) cleanByTime() {
	for {
		<-time.After(c.maxWait)

		log.Debugf("cleanup by time complite")
		c.c <- c.items
		c.items = make([]interface{}, 0)
	}
}

func (c *Cache) cleanBySize() {
	for {
		if len(c.items) > c.maxSize {
			log.Debugf("cleanup by size complite")
			c.c <- c.items
			c.items = make([]interface{}, 0)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Cache) cleanByteSize() {
	for {
		size := cap(c.items)
		if size > c.maxByteSize {
			log.Debugf("cleanup by byte size complite")
			c.c <- c.items
			c.items = make([]interface{}, 0)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
