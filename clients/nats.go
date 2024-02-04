package clients

import (
	"fmt"
	"github.com/Kebastos/NatsToCh/config"
	"github.com/Kebastos/NatsToCh/log"
	"github.com/Kebastos/NatsToCh/metrics"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

type NatsClient struct {
	mx  sync.Mutex
	cfg *config.NATSConfig
	nc  *nats.Conn
}

func NewNatsClient(cfg *config.NATSConfig) *NatsClient {
	return &NatsClient{
		cfg: cfg,
	}
}

func (c *NatsClient) Connect() error {
	c.mx.Lock()
	defer c.mx.Unlock()

	options := []nats.Option{
		nats.Name(c.cfg.ClientName),
		nats.MaxReconnects(c.cfg.MaxReconnect),
		nats.ReconnectWait(time.Duration(c.cfg.ReconnectWait) * time.Millisecond),
		nats.Timeout(time.Duration(c.cfg.ConnectTimeout) * time.Millisecond),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Errorf("NATS error: %s\n", err)
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Errorf("NATS disconnected: %s\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Errorf("NATS connection closed")
		}),
		nats.ConnectHandler(func(nc *nats.Conn) {
			log.Infof("NATS connected")
		}),
		nats.NoCallbacksAfterClientClose(),
	}

	nc, err := nats.Connect(c.cfg.Server, options...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	c.nc = nc
	return nil
}

func (c *NatsClient) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.nc == nil {
		return nil, fmt.Errorf("nats connection is not available")
	}

	return c.nc.Subscribe(subject, func(msg *nats.Msg) {
		metrics.GotMessageCount.Inc()
		handler(msg)
	})
}

func (c *NatsClient) QueueSubscribe(subject string, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.nc == nil {
		return nil, fmt.Errorf("nats connection is not available")
	}

	return c.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		metrics.GotMessageCount.Inc()
		handler(msg)
	})
}
