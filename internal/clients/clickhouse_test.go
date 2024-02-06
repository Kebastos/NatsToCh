package clients_test

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/models"
	"os"
	"testing"
	"time"
)

const (
	tableName = "test"
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
)

func TestMain(m *testing.M) {
	metrics.MustRegister()

	code := m.Run()

	os.Exit(code)
}

func TestNewClickhouseClientConnect(t *testing.T) {
	client := clients.NewClickhouseClient(cfg)

	err := client.Connect()
	if err != nil {
		t.Errorf("Failed to connect to Clickhouse server: %s", err)
	}
}

func TestNewClickhouseClientBatchInsertToDefaultSchema(t *testing.T) {
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsertToDefaultSchema(context.Background(), tableName, []interface{}{testData})
	if err != nil {
		t.Errorf("Failed to batch insert to default schema: %s", err)
	}
}

func TestNewClickhouseClientAsyncInsertToDefaultSchema(t *testing.T) {
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultSchema(context.Background(), tableName, []interface{}{testData}, true)
	if err != nil {
		t.Errorf("Failed to async insert to default schema: %s", err)
	}
}

func TestNewClickhouseClientAsyncInsertToDefaultSchemaNoData(t *testing.T) {
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultSchema(context.Background(), tableName, []interface{}{}, true)
	if err == nil {
		t.Errorf("Expected error when async inserting with no data")
	}
}
