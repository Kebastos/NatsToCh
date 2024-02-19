package workers_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/workers"
	"github.com/nats-io/nats.go"
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

func TestNatsWorkerStartsWithBuffer(t *testing.T) {
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
	worker := workers.NewNatsWorker(cfg, sb, ch, logger)

	err := worker.Start(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestNatsWorkerStartsWithNoBuffer(t *testing.T) {
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
	worker := workers.NewNatsWorker(cfg, sb, ch, logger)

	err := worker.Start(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestNatsWorkerStartsWithNoBufferAsync(t *testing.T) {
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
	worker := workers.NewNatsWorker(cfg, sb, ch, logger)

	err := worker.Start(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestNatsWorkerStartsWithError(t *testing.T) {
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
	worker := workers.NewNatsWorker(cfg, sb, ch, logger)

	err := worker.Start(context.Background())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
