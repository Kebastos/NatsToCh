package clients_test

import (
	"context"
	"testing"

	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
)

func TestNewClickhouseClientConnect(t *testing.T) {
	cfg := &config.CHConfig{
		Host:            "localhost",
		Port:            9000,
		User:            "default",
		Password:        "",
		Database:        "default",
		ConnMaxLifetime: 0,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	client := clients.NewClickhouseClient(cfg)

	err := client.Connect()
	if err != nil {
		t.Errorf("Failed to connect to Clickhouse server: %s", err)
	}
}

func TestNewClickhouseClientBatchInsertToDefaultSchema(t *testing.T) {
	cfg := &config.CHConfig{
		Host:            "localhost",
		Port:            9000,
		User:            "default",
		Password:        "",
		Database:        "default",
		ConnMaxLifetime: 0,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.BatchInsertToDefaultSchema(context.Background(), "test_table", []interface{}{"test_data"})
	if err != nil {
		t.Errorf("Failed to batch insert to default schema: %s", err)
	}
}

func TestNewClickhouseClientAsyncInsertToDefaultSchema(t *testing.T) {
	cfg := &config.CHConfig{
		Host:            "localhost",
		Port:            9000,
		User:            "default",
		Password:        "",
		Database:        "default",
		ConnMaxLifetime: 0,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultSchema(context.Background(), "test_table", []interface{}{"test_data"}, true)
	if err != nil {
		t.Errorf("Failed to async insert to default schema: %s", err)
	}
}

func TestNewClickhouseClientAsyncInsertToDefaultSchemaNoData(t *testing.T) {
	cfg := &config.CHConfig{
		Host:            "localhost",
		Port:            9000,
		User:            "default",
		Password:        "",
		Database:        "default",
		ConnMaxLifetime: 0,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	}
	client := clients.NewClickhouseClient(cfg)
	err := client.Connect()
	if err != nil {
		t.Errorf("failed to connect to Clickhouse server: %s", err)
	}

	err = client.AsyncInsertToDefaultSchema(context.Background(), "test_table", []interface{}{}, true)
	if err == nil {
		t.Errorf("Expected error when async inserting with no data")
	}
}
