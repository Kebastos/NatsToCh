package workers

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/cache"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/models"
	"github.com/nats-io/nats.go"
	"time"
)

type ClickhouseStorage interface {
	BatchInsertToDefaultSchema(ctx context.Context, tableName string, data []interface{}) error
	AsyncInsertToDefaultSchema(ctx context.Context, tableName string, data []interface{}, wait bool) error
}

type NatsSub interface {
	Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
	QueueSubscribe(subject string, queue string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
}

type NatsWorker struct {
	cfg    *config.Config
	sb     NatsSub
	ch     ClickhouseStorage
	logger *log.Log
}

func NewNatsWorker(cfg *config.Config, sb NatsSub, ch ClickhouseStorage, logger *log.Log) *NatsWorker {
	return &NatsWorker{
		cfg:    cfg,
		sb:     sb,
		ch:     ch,
		logger: logger,
	}
}

func (n *NatsWorker) Start(ctx context.Context) error {
	for _, s := range n.cfg.Subjects {
		var c func(m *nats.Msg)
		switch {
		case s.UseBuffer:
			c = n.callbackWithBuffer(ctx, s)
		case s.Async:
			c = n.callbackNoBufferAsync(ctx, s.TableName, s.AsyncInsertConfig.Wait)
		default:
			c = n.callbackNoBuffer(ctx, s.TableName)
		}

		err := n.subs(s, c)
		if err != nil {
			return err
		}

		n.logger.Infof("subscribed to %s\n with params: buffer=%t and table=%s", s.Name, s.UseBuffer, s.TableName)
	}

	return nil
}

func (n *NatsWorker) subs(cfg config.Subject, f func(m *nats.Msg)) error {
	var err error
	if cfg.Queue != "" {
		_, err = n.sb.QueueSubscribe(cfg.Name, cfg.Queue, f)
	} else {
		_, err = n.sb.Subscribe(cfg.Name, f)
	}

	return err
}

func (n *NatsWorker) callbackWithBuffer(ctx context.Context, cfg config.Subject) func(m *nats.Msg) {
	c := make(chan []interface{}, 1)
	cc := cache.New(&cfg.BufferConfig, n.logger, c)
	cw := NewClickhouseWorker(&cfg, n.ch, c, n.logger)
	cw.Start(ctx)

	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		cc.Set(entity)
	}

	return callback
}

func (n *NatsWorker) callbackNoBuffer(ctx context.Context, table string) func(m *nats.Msg) {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.BatchInsertToDefaultSchema(ctx, table, []interface{}{entity})
		if err != nil {
			n.logger.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	return callback
}

func (n *NatsWorker) callbackNoBufferAsync(ctx context.Context, table string, wait bool) func(m *nats.Msg) {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.AsyncInsertToDefaultSchema(ctx, table, []interface{}{entity}, wait)
		if err != nil {
			n.logger.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	return callback
}
