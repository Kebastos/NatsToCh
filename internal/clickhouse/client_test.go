package clickhouse

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/models"
	"testing"
)

const (
	tableName  = "test"
	wrongTable = "wrongTable"
)

var (
	testData = models.NewDefaultEntity("test_client", "test_subject", "test_data")
	cfg      = &config.CHConfig{
		Host:            "localhost",
		Port:            9000,
		User:            "default",
		Password:        "",
		Database:        "test",
		ConnMaxLifetime: 0,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	logger = log.MustConfig()
)

type MockMetrics struct{}

func (m *MockMetrics) InsertMessageCountAdd(_ string, _ int) {}

func TestClickhouseClientConnect(t *testing.T) {
	client := NewClickhouseClient(cfg, logger, &MockMetrics{})

	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}
}

func TestClickhouseClientConnectIsClosed(t *testing.T) {
	client := NewClickhouseClient(cfg, logger, &MockMetrics{})

	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("failed to close Clickhouse client: %s", err)
	}
	stats := client.ConnStatus()
	if stats.Open > 0 {
		t.Errorf("expected 0 open connections, got %d", stats.Open)
	}
}

func TestClickhouseClientBatchInsertToDefaultTable(t *testing.T) {
	client := NewClickhouseClient(cfg, logger, &MockMetrics{})
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsert(context.Background(), tableName, []interface{}{testData})
	if err != nil {
		t.Errorf("failed to batch insert to default schema: %s", err)
	}
}

func TestClickhouseClientBatchInsertToDefaultTableNoData(t *testing.T) {
	client := NewClickhouseClient(cfg, logger, &MockMetrics{})
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsert(context.Background(), tableName, []interface{}{})
	if err == nil {
		t.Errorf("expected error when async inserting with no data")
	}
}

func TestClickhouseClientBatchInsertToWrongTable(t *testing.T) {
	client := NewClickhouseClient(cfg, logger, &MockMetrics{})
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsert(context.Background(), wrongTable, []interface{}{testData})
	if err == nil {
		t.Errorf("expected error when async inserting with")
	}
}
