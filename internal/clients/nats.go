package clients

import (
	"fmt"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

type NatsClient struct {
	mx     sync.Mutex
	cfg    *config.NATSConfig
	nc     *nats.Conn
	logger *log.Log
}

func NewNatsClient(cfg *config.NATSConfig, logger *log.Log) *NatsClient {
	return &NatsClient{
		cfg:    cfg,
		logger: logger,
	}
}

func (c *NatsClient) Connect() error {
	c.mx.Lock()
	defer c.mx.Unlock()

	options := []nats.Option{
		nats.Name(c.cfg.ClientName),
		nats.UserInfo(c.cfg.User, c.cfg.Password),
		nats.MaxReconnects(c.cfg.MaxReconnect),
		nats.ReconnectWait(time.Duration(c.cfg.ReconnectWait) * time.Second),
		nats.Timeout(time.Duration(c.cfg.ConnectTimeout) * time.Second),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			c.logger.Errorf("NATS error: %s\n", err)
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			c.logger.Errorf("NATS disconnected: %s\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			c.logger.Infof("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			c.logger.Errorf("NATS connection closed")
		}),
		nats.ConnectHandler(func(nc *nats.Conn) {
			c.logger.Infof("connected to Nats server at %s", c.cfg.Server)
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

func (c *NatsClient) Shutdown() {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.nc != nil {
		err := c.nc.Drain()
		if err != nil {
			c.logger.Fatalf("failed to drain NATS connection: %s", err)
		}
		c.nc.Close()
	}

	c.logger.Infof("nats client was disconnected")
}

func (c *NatsClient) ConnStatus() nats.Status {
	if c.nc != nil {
		return c.nc.Status()
	}
	return nats.DISCONNECTED
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
