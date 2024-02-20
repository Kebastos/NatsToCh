package cache

import (
	"github.com/Kebastos/NatsToCh/internal/log"
	"os"
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/internal/config"
)

type MockMetrics struct{}

func (m *MockMetrics) QueueMessageCountInc(_ string)   {}
func (m *MockMetrics) QueueMessageCountDrain(_ string) {}

var (
	logger   = log.MustConfig()
	c        = make(chan []interface{})
	shortCfg = config.Subject{

		Name:      "test",
		UseBuffer: true,
		BufferConfig: config.BufferConfig{
			MaxSize: 10,
			MaxWait: 1 * time.Second,
		},
	}
	longCfg = config.Subject{

		Name:      "test",
		UseBuffer: true,
		BufferConfig: config.BufferConfig{
			MaxSize: 10,
			MaxWait: 600 * time.Second,
		},
	}
)

func TestMain(m *testing.M) {
	go func() {
		for {
			<-c
			logger.Infof("got message from cache")
		}
	}()

	code := m.Run()

	os.Exit(code)
}

func TestNewCache(t *testing.T) {
	ch := New(&shortCfg, logger, c, &MockMetrics{})

	if ch.Count() != 0 {
		t.Errorf("new cache should be empty, got %d items", ch.Count())
	}
}

func TestCacheSet(t *testing.T) {
	ch := New(&shortCfg, logger, c, &MockMetrics{})

	ch.Set("test")

	if ch.Count() != 1 {
		t.Errorf("cache count should be 1, got %d", ch.Count())
	}
}

func TestCacheDrainAtLenOverflow(t *testing.T) {
	ch := New(&longCfg, logger, c, &MockMetrics{})
	ch.StartCleaner()

	ch.Set("test1")
	ch.Set("test2")
	ch.Set("test3")

	<-time.After(1 * time.Second)

	if ch.Count() != 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}

func TestCacheCleanByTime(t *testing.T) {
	ch := New(&shortCfg, logger, c, &MockMetrics{})

	ch.StartCleaner()

	ch.Set("test")

	<-time.After(2 * time.Second)
	if ch.Count() > 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}

func TestCacheShutdown(t *testing.T) {
	ch := New(&longCfg, logger, c, &MockMetrics{})
	ch.StartCleaner()

	ch.Set("test1")
	ch.Set("test2")
	ch.Set("test3")

	ch.Shutdown()
	<-time.After(1 * time.Second)
	if ch.Count() != 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}

func TestCacheCloseByShutdown(t *testing.T) {
	ch := New(&longCfg, logger, c, &MockMetrics{})
	ch.StartCleaner()

	ch.Shutdown()
	ch.Set("test1")

	<-time.After(1 * time.Second)

	if ch.Count() != 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}
