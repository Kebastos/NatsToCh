package nats_test

import (
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	client "github.com/Kebastos/NatsToCh/internal/nats"
	"github.com/nats-io/nats.go"
	"testing"
)

var (
	natsCfg = &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	logger = log.MustConfig()
)

type MockMetrics struct{}

func (m *MockMetrics) GotMessageCountInc(_ string) {}

func TestNatsClientConnect(t *testing.T) {
	c := client.NewClient(natsCfg, logger, &MockMetrics{})

	err := c.Connect()
	if err != nil {
		t.Errorf("Failed to connect to NATS server: %s", err)
	}
}

func TestNatsClientShutdown(t *testing.T) {
	c := client.NewClient(natsCfg, logger, &MockMetrics{})

	err := c.Connect()
	if err != nil {
		t.Errorf("Failed to connect to NATS server: %s", err)
	}
	c.Shutdown()
	st := c.ConnStatus()
	if st != nats.CLOSED {
		t.Errorf("Expected connection to be closed, got %s", st.String())
	}
}

func TestNatsClientSubscribe(t *testing.T) {
	c := client.NewClient(natsCfg, logger, &MockMetrics{})
	err := c.Connect()
	if err != nil {
		t.Errorf("failed to connect to NATS server: %s", err)
	}

	_, err = c.Subscribe("test.subject", func(msg *nats.Msg) {})
	if err != nil {
		t.Errorf("failed to subscribe to subject: %s", err)
	}
}

func TestNatsClientQueueSubscribe(t *testing.T) {
	c := client.NewClient(natsCfg, logger, &MockMetrics{})
	err := c.Connect()
	if err != nil {
		t.Errorf("failed to connect to NATS server: %s", err)
	}

	_, err = c.QueueSubscribe("test.subject", "test.queue", func(msg *nats.Msg) {})
	if err != nil {
		t.Errorf("failed to queue subscribe to subject: %s", err)
	}
}

func TestNatsClientSubscribeWithoutConnect(t *testing.T) {
	cfg := &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	c := client.NewClient(cfg, logger, &MockMetrics{})

	_, err := c.Subscribe("test.subject", func(msg *nats.Msg) {})
	if err == nil {
		t.Errorf("expected error when subscribing without connecting")
	}
}

func TestNatsClientQueueSubscribeWithoutConnect(t *testing.T) {
	c := client.NewClient(natsCfg, logger, &MockMetrics{})

	_, err := c.QueueSubscribe("test.subject", "test.queue", func(msg *nats.Msg) {})
	if err == nil {
		t.Errorf("expected error when queue subscribing without connecting")
	}
}
