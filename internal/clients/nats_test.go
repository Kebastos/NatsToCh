package clients_test

import (
	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/nats-io/nats.go"
	"testing"
)

func TestNatsClientConnect(t *testing.T) {
	cfg := &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	client := clients.NewNatsClient(cfg)

	err := client.Connect()
	if err != nil {
		t.Errorf("Failed to connect to NATS server: %s", err)
	}
}

func TestNatsClientSubscribe(t *testing.T) {
	cfg := &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	client := clients.NewNatsClient(cfg)
	client.Connect()

	_, err := client.Subscribe("test.subject", func(msg *nats.Msg) {})
	if err != nil {
		t.Errorf("Failed to subscribe to subject: %s", err)
	}
}

func TestNatsClientQueueSubscribe(t *testing.T) {
	cfg := &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	client := clients.NewNatsClient(cfg)
	client.Connect()

	_, err := client.QueueSubscribe("test.subject", "test.queue", func(msg *nats.Msg) {})
	if err != nil {
		t.Errorf("Failed to queue subscribe to subject: %s", err)
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
	client := clients.NewNatsClient(cfg)

	_, err := client.Subscribe("test.subject", func(msg *nats.Msg) {})
	if err == nil {
		t.Errorf("Expected error when subscribing without connecting")
	}
}

func TestNatsClientQueueSubscribeWithoutConnect(t *testing.T) {
	cfg := &config.NATSConfig{
		Server:         "nats://localhost:4222",
		ClientName:     "TestClient",
		MaxReconnect:   10,
		ReconnectWait:  500,
		ConnectTimeout: 2000,
	}
	client := clients.NewNatsClient(cfg)

	_, err := client.QueueSubscribe("test.subject", "test.queue", func(msg *nats.Msg) {})
	if err == nil {
		t.Errorf("Expected error when queue subscribing without connecting")
	}
}
