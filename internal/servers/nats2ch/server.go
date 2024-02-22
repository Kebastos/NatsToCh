package nats2ch

import (
	"context"
	"fmt"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/nats-io/nats.go"
)

type Instrumentation interface {
	QueueMessageCountInc(name string)
	QueueMessageCountDrain(name string)
}

type ClickhouseStorage interface {
	BatchInsert(ctx context.Context, tableName string, data []interface{}) error
}

type NatsSub interface {
	Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
	QueueSubscribe(subject string, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
}

type Nats2Ch struct {
	cfg     *config.Config
	sb      NatsSub
	ch      ClickhouseStorage
	logger  *log.Log
	metrics Instrumentation
}

func NewServer(cfg *config.Config, sb NatsSub, ch ClickhouseStorage, logger *log.Log, metrics Instrumentation) *Nats2Ch {
	return &Nats2Ch{
		cfg:     cfg,
		sb:      sb,
		ch:      ch,
		logger:  logger,
		metrics: metrics,
	}
}

func (n *Nats2Ch) Start(ctx context.Context) error {
	for _, s := range n.cfg.Subjects {
		var c func(m *nats.Msg)
		switch {
		case s.UseBuffer:
			if s.BufferConfig.MaxSize == 0 || s.BufferConfig.MaxWait == 0 {
				return fmt.Errorf("buffer configuration must be greater 0  for %s", s.Name)
			}
			c = n.callbackWithBuffer(ctx, s)
		default:
			c = n.callbackNoBuffer(ctx, s.TableName)
		}

		err := n.subs(s, c)
		if err != nil {
			return err
		}

		n.logger.Infof("subscribed to %s with params: buffer=%t and table=%s", s.Name, s.UseBuffer, s.TableName)
	}

	return nil
}

func (n *Nats2Ch) subs(cfg config.Subject, f func(m *nats.Msg)) error {
	var err error
	if cfg.Queue != "" {
		_, err = n.sb.QueueSubscribe(cfg.Name, cfg.Queue, f)
	} else {
		_, err = n.sb.Subscribe(cfg.Name, f)
	}

	return err
}

func (n *Nats2Ch) startInsert(ctx context.Context, c chan []interface{}, tableName string) {
	go func() {
		for {
			items := <-c
			if len(items) > 0 {
				err := n.ch.BatchInsert(ctx, tableName, items)
				if err != nil {
					n.logger.Errorf("failed to insert data to clickhouse: %s", err)
				}
			}
		}
	}()
}
