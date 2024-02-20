package cache

import (
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"sync"
	"time"
)

type Instrumentation interface {
	QueueMessageCountInc(name string)
	QueueMessageCountDrain(name string)
}

type Cache struct {
	sync.RWMutex
	logger  *log.Log
	maxSize int
	maxWait time.Duration
	items   []interface{}
	cfg     *config.Subject
	c       chan []interface{}
	stop    bool
	metrics Instrumentation
}

func New(cfg *config.Subject, logger *log.Log, c chan []interface{}, metrics Instrumentation) *Cache {
	items := make([]interface{}, 0)

	cache := Cache{
		items:   items,
		logger:  logger,
		maxSize: cfg.BufferConfig.MaxSize,
		maxWait: cfg.BufferConfig.MaxWait,
		cfg:     cfg,
		c:       c,
		metrics: metrics,
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
	c.metrics.QueueMessageCountInc(c.cfg.Name)
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
	c.metrics.QueueMessageCountDrain(c.cfg.Name)
}

func (c *Cache) drainByTimeout() {
	for {
		c.drain("timeout")
		<-time.After(c.maxWait)
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
