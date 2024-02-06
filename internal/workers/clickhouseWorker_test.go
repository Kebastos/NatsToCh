package workers_test

import (
	"context"
	"testing"

	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/workers"
)

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

var cfg = &config.Subject{
	TableName: "test_table",
}

func TestClickhouseWorkerStart(t *testing.T) {
	ch := &MockClickhouseStorage{
		BatchInsertToDefaultSchemaFunc: func(ctx context.Context, tableName string, items []interface{}) error {
			return nil
		},
	}
	c := make(chan []interface{}, 1)
	worker := workers.NewClickhouseWorker(cfg, ch, c)

	worker.Start(context.Background())
	c <- []interface{}{"test_data"}

	<-c
}
