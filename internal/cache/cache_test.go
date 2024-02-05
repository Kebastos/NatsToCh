package cache_test

import (
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/internal/cache"
	"github.com/Kebastos/NatsToCh/internal/config"
)

func TestNewCache(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize:     10,
		MaxWait:     1 * time.Second,
		MaxByteSize: 100,
	}
	c := make(chan []interface{})
	ch := cache.New(cfg, c)

	if ch.Count() != 0 {
		t.Errorf("New cache should be empty, got %d items", ch.Count())
	}
}

func TestCacheSet(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize:     10,
		MaxWait:     1 * time.Second,
		MaxByteSize: 100,
	}
	c := make(chan []interface{})
	ch := cache.New(cfg, c)

	ch.Set("test")

	if ch.Count() != 1 {
		t.Errorf("Cache count should be 1, got %d", ch.Count())
	}
}

func TestCacheOverflow(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize:     2,
		MaxWait:     1 * time.Second,
		MaxByteSize: 100,
	}
	c := make(chan []interface{})
	ch := cache.New(cfg, c)

	ch.Set("test1")
	ch.Set("test2")
	ch.Set("test3")

	select {
	case <-c:
		if ch.Count() != 1 {
			t.Errorf("Cache count should be 1, got %d", ch.Count())
		}
	case <-time.After(2 * time.Second):
		t.Errorf("Cache did not overflow in time")
	}
}

func TestCacheCleanByTime(t *testing.T) {
	cfg := &config.BufferConfig{
		MaxSize:     10,
		MaxWait:     1 * time.Second,
		MaxByteSize: 100,
	}
	c := make(chan []interface{})
	ch := cache.New(cfg, c)

	ch.Set("test")

	select {
	case <-c:
		if ch.Count() != 0 {
			t.Errorf("Cache count should be 0, got %d", ch.Count())
		}
	case <-time.After(2 * time.Second):
		t.Errorf("Cache did not clean in time")
	}
}
