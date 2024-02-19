package cache

import (
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	logger  *log.Log
	maxSize int
	maxWait time.Duration
	items   []interface{}
	cfg     *config.BufferConfig
	c       chan []interface{}
	stop    bool
}

func New(cfg *config.BufferConfig, logger *log.Log, c chan []interface{}) *Cache {
	items := make([]interface{}, 0)

	cache := Cache{
		items:   items,
		logger:  logger,
		maxSize: cfg.MaxSize,
		maxWait: cfg.MaxWait,
		cfg:     cfg,
		c:       c,
	}

	return &cache
}

func (c *Cache) StartCleaner() {
	go c.drainByTimeout()
	go c.drainByLen()
}

func (c *Cache) Shutdown() {
	c.stop = true
	c.drain("shutdown")
}

func (c *Cache) Set(value interface{}) {
	c.Lock()
	defer c.Unlock()

	if c.stop {
		c.logger.Errorf("cache closed")
		return
	}
	c.items = append(c.items, value)
}

func (c *Cache) Count() int {
	return len(c.items)
}

func (c *Cache) drain(cleanType string) {
	c.Lock()
	defer c.Unlock()

	c.logger.Debugf("cleanup by %s complete", cleanType)
	c.c <- c.items
	c.items = make([]interface{}, 0)
}

func (c *Cache) drainByTimeout() {
	for {
		<-time.After(c.maxWait)
		c.drain("timeout")
	}
}

func (c *Cache) drainByLen() {
	for {
		if len(c.items) > c.maxSize {
			c.drain("length")
		}
		<-time.After(5 * time.Second)
	}
}
