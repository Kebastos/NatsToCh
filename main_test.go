package main

import (
	"context"
	"testing"
	"time"

	"github.com/Kebastos/NatsToCh/clickhouse"
	"github.com/Kebastos/NatsToCh/config"
	"github.com/Kebastos/NatsToCh/workers"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNatsClient struct {
	mock.Mock
}

func (m *MockNatsClient) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNatsClient) Subscribe(subject string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subject, cb)
	return nil, args.Error(1)
}

type MockClickhouseClient struct {
	mock.Mock
}

func (m *MockClickhouseClient) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClickhouseClient) BatchInsertToDefaultSchema(tableName string, data []interface{}, ctx context.Context) error {
	args := m.Called(tableName, data, ctx)
	return args.Error(0)
}

func TestMainFunctionStartsSuccessfully(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Debug: true,
			HTTP: config.HTTP{
				ListenAddr: "localhost:8080",
			},
		},
		NATSConfig: nats.Config{
			URL: "nats://localhost:4222",
		},
		CHConfig: clickhouse.Config{
			URL: "tcp://localhost:9000?debug=true",
		},
	}

	mockNatsClient := new(MockNatsClient)
	mockNatsClient.On("Connect").Return(nil)

	mockClickhouseClient := new(MockClickhouseClient)
	mockClickhouseClient.On("Connect").Return(nil)

	natsWorker := workers.NewNatsWorker(cfg, mockNatsClient, mockClickhouseClient)
	natsWorker.On("Start").Return(nil)

	go main()

	time.Sleep(time.Second) // give some time for the main function to run

	mockNatsClient.AssertExpectations(t)
	mockClickhouseClient.AssertExpectations(t)
}

func TestMainFunctionFailsToStart(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Debug: true,
			HTTP: config.HTTP{
				ListenAddr: "localhost:8080",
			},
		},
		NATSConfig: nats.Config{
			URL: "nats://localhost:4222",
		},
		CHConfig: clickhouse.Config{
			URL: "tcp://localhost:9000?debug=true",
		},
	}

	mockNatsClient := new(MockNatsClient)
	mockNatsClient.On("Connect").Return(assert.AnError)

	mockClickhouseClient := new(MockClickhouseClient)
	mockClickhouseClient.On("Connect").Return(nil)

	natsWorker := workers.NewNatsWorker(cfg, mockNatsClient, mockClickhouseClient)
	natsWorker.On("Start").Return(nil)

	go main()

	time.Sleep(time.Second) // give some time for the main function to run

	mockNatsClient.AssertExpectations(t)
	mockClickhouseClient.AssertExpectations(t)
}
