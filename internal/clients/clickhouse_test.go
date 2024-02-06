package clients_test

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/models"
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
)

const (
	tableName = "test"
)

var testData = []interface{}{
	&models.DefaultTable{
		Subject:        "test_subject",
		CreateDateTime: time.Now(),
		Content:        "test_data",
	},
}

var cfg = &config.CHConfig{
	Host:            "localhost",
	Port:            9000,
	User:            "default",
	Password:        "",
	Database:        "test",
	ConnMaxLifetime: 0,
	MaxOpenConns:    10,
	MaxIdleConns:    5,
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

	err = client.BatchInsertToDefaultSchema(context.Background(), tableName, testData)
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

	err = client.AsyncInsertToDefaultSchema(context.Background(), tableName, testData, true)
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
