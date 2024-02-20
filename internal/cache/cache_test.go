package cache

import (
	"github.com/Kebastos/NatsToCh/internal/log"
	"os"
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/internal/config"
)

var (
	logger = log.MustConfig()
	c      = make(chan []interface{})
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
	cfg := &config.BufferConfig{
		MaxSize: 10,
		MaxWait: 1 * time.Second,
	}

	ch := New(cfg, logger, c)

	if ch.Count() != 0 {
		t.Errorf("new cache should be empty, got %d items", ch.Count())
	}
}

func TestCacheSet(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize: 10,
		MaxWait: 1 * time.Second,
	}

	ch := New(cfg, logger, c)

	ch.Set("test")

	if ch.Count() != 1 {
		t.Errorf("cache count should be 1, got %d", ch.Count())
	}
}

func TestCacheDrainAtLenOverflow(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize: 2,
		MaxWait: 600 * time.Second,
	}

	ch := New(cfg, logger, c)
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
	cfg := &config.BufferConfig{
		MaxSize: 100,
		MaxWait: 1 * time.Second,
	}

	ch := New(cfg, logger, c)

	ch.StartCleaner()

	ch.Set("test")

	<-time.After(2 * time.Second)
	if ch.Count() > 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}

func TestCacheShutdown(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize: 10,
		MaxWait: 600 * time.Second,
	}
	ch := New(cfg, logger, c)
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
	cfg := &config.BufferConfig{
		MaxSize: 10,
		MaxWait: 600 * time.Second,
	}
	ch := New(cfg, logger, c)
	ch.StartCleaner()

	ch.Shutdown()
	ch.Set("test1")

	<-time.After(1 * time.Second)

	if ch.Count() != 0 {
		t.Errorf("cache count should be 0, got %d", ch.Count())
	}
}
