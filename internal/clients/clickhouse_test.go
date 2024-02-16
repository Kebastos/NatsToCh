package clients_test

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/models"
	"os"
	"testing"
	"time"
)

const (
	tableName  = "test"
	wrongTable = "wrongTable"
)

var (
	testData = &models.DefaultTable{
		Subject:        "test_subject",
		CreateDateTime: time.Now(),
		Content:        "test_data",
	}
	cfg = &config.CHConfig{
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

func TestMain(m *testing.M) {
	metrics.MustRegister()

	code := m.Run()

	os.Exit(code)
}

func TestClickhouseClientConnect(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)

	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}
}

func TestClickhouseClientConnectIsClosed(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)

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
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsertToDefaultTable(context.Background(), tableName, []interface{}{testData})
	if err != nil {
		t.Errorf("failed to batch insert to default schema: %s", err)
	}
}

func TestClickhouseClientBatchInsertToDefaultTableNoData(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsertToDefaultTable(context.Background(), tableName, []interface{}{})
	if err == nil {
		t.Errorf("expected error when async inserting with no data")
	}
}

func TestClickhouseClientBatchInsertToWrongTable(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsertToDefaultTable(context.Background(), wrongTable, []interface{}{testData})
	if err == nil {
		t.Errorf("expected error when async inserting with")
	}
}

func TestClickhouseClientAsyncInsertToDefaultTable(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultTable(context.Background(), tableName, []interface{}{testData}, true)
	if err != nil {
		t.Errorf("failed to async insert to default schema: %s", err)
	}
}

func TestClickhouseClientAsyncInsertToDefaultSchemaNoData(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultTable(context.Background(), tableName, []interface{}{}, true)
	if err == nil {
		t.Errorf("expected error when async inserting with no data")
	}
}

func TestClickhouseClientAsyncInsertToDefaultWrongTable(t *testing.T) {
	client := clients.NewClickhouseClient(cfg, logger)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultTable(context.Background(), wrongTable, []interface{}{testData}, true)
	if err == nil {
		t.Errorf("expected error when async inserting with")
	}
}
