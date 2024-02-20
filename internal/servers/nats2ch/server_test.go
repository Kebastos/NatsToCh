package nats2ch

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

type MockNatsSub struct {
	SubscribeFunc      func(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
	QueueSubscribeFunc func(subject string, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
}

func (m *MockNatsSub) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return m.SubscribeFunc(subject, handler)
}

func (m *MockNatsSub) QueueSubscribe(subject string, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return m.QueueSubscribeFunc(subject, queue, handler)
}

type MockClickhouseStorage struct {
	BatchInsertToDefaultSchemaFunc func(ctx context.Context, tableName string, items []interface{}) error
	AsyncInsertToDefaultSchemaFunc func(ctx context.Context, tableName string, data []interface{}, wait bool) error
}

func (m *MockClickhouseStorage) BatchInsertToDefaultSchema(ctx context.Context, tableName string, items []interface{}) error {
	return m.BatchInsertToDefaultSchemaFunc(ctx, tableName, items)
}

func (m *MockClickhouseStorage) AsyncInsertToDefaultSchema(ctx context.Context, tableName string, data []interface{}, wait bool) error {
	return m.AsyncInsertToDefaultSchemaFunc(ctx, tableName, data, wait)
}

type MockMetrics struct{}

func (m *MockMetrics) QueueMessageCountInc(_ string)   {}
func (m *MockMetrics) QueueMessageCountDrain(_ string) {}

var logger = log.MustConfig()

func TestNats2Ch_StartWithBuffer(t *testing.T) {
	cfg := &config.Config{
		Subjects: []config.Subject{
			{
				Name:      "test_subject",
				TableName: "test_table",
				UseBuffer: true,
				BufferConfig: config.BufferConfig{
					MaxSize: 10,
					MaxWait: 10 * time.Second,
				},
			},
		},
	}
	ch := &MockClickhouseStorage{
		BatchInsertToDefaultSchemaFunc: func(ctx context.Context, tableName string, items []interface{}) error {
			return nil
		},
	}
	sb := &MockNatsSub{
		SubscribeFunc: func(subject string, cb func(msg *nats.Msg)) (*nats.Subscription, error) {
			return nil, nil
		},
	}
	srv := NewServer(cfg, sb, ch, logger, &MockMetrics{})

	err := srv.Start(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestNats2Ch_StartWithNoBuffer(t *testing.T) {
	cfg := &config.Config{
		Subjects: []config.Subject{
			{
				Name:      "test_subject",
				TableName: "test_table",
			},
		},
	}
	ch := &MockClickhouseStorage{
		BatchInsertToDefaultSchemaFunc: func(ctx context.Context, tableName string, items []interface{}) error {
			return nil
		},
	}
	sb := &MockNatsSub{
		SubscribeFunc: func(subject string, cb func(msg *nats.Msg)) (*nats.Subscription, error) {
			return nil, nil
		},
	}
	srv := NewServer(cfg, sb, ch, logger, &MockMetrics{})

	err := srv.Start(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestNats2Ch_StartWithNoBufferAsync(t *testing.T) {
	cfg := &config.Config{
		Subjects: []config.Subject{
			{
				Name:      "test_subject",
				TableName: "test_table",
				Async:     true,
				AsyncInsertConfig: config.AsyncInsertConfig{
					Wait: true,
				},
			},
		},
	}
	ch := &MockClickhouseStorage{
		AsyncInsertToDefaultSchemaFunc: func(ctx context.Context, tableName string, data []interface{}, wait bool) error {
			return nil
		},
	}
	sb := &MockNatsSub{
		SubscribeFunc: func(subject string, cb func(msg *nats.Msg)) (*nats.Subscription, error) {
			return nil, nil
		},
	}
	srv := NewServer(cfg, sb, ch, logger, &MockMetrics{})

	err := srv.Start(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestNats2Ch_StartWithError(t *testing.T) {
	cfg := &config.Config{
		Subjects: []config.Subject{
			{
				Name:      "test_subject",
				TableName: "test_table",
			},
		},
	}
	ch := &MockClickhouseStorage{
		BatchInsertToDefaultSchemaFunc: func(ctx context.Context, tableName string, items []interface{}) error {
			return errors.New("batch insert error")
		},
	}
	sb := &MockNatsSub{
		SubscribeFunc: func(subject string, cb func(msg *nats.Msg)) (*nats.Subscription, error) {
			return nil, fmt.Errorf("subscribe error")
		},
	}
	srv := NewServer(cfg, sb, ch, logger, &MockMetrics{})

	err := srv.Start(context.Background())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}