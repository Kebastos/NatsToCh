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
}

type NatsWorker struct {
	cfg *config.Config
	sb  NatsSub
	ch  ClickhouseStorage
}

func NewNatsWorker(cfg *config.Config, sb NatsSub, ch ClickhouseStorage) *NatsWorker {
	return &NatsWorker{
		cfg: cfg,
		sb:  sb,
		ch:  ch,
	}
}

func (n *NatsWorker) Start(ctx context.Context) error {
	for _, s := range n.cfg.Subjects {
		var err error

		switch {
		case s.UseBuffer:
			err = n.subsWithBuffer(ctx, s.Name, s)
		case s.Async:
			err = n.subsNoBufferAsync(ctx, s.Name, s.TableName, s.AsyncInsertConfig.Wait)
		default:
			err = n.subsNoBuffer(ctx, s.Name, s.TableName)
		}

		if err != nil {
			return err
		}

		log.Infof("subscribed to %s\n with params: buffer=%t and table=%s", s.Name, s.UseBuffer, s.TableName)
	}

	return nil
}

func (n *NatsWorker) subsWithBuffer(ctx context.Context, subject string, cfg config.Subject) error {
	c := make(chan []interface{}, 1)
	cc := cache.New(&cfg.BufferConfig, c)
	cw := NewClickhouseWorker(&cfg, n.ch, c)
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

	_, err := n.sb.Subscribe(subject, callback)

	return err
}

func (n *NatsWorker) subsNoBuffer(ctx context.Context, subject string, table string) error {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.BatchInsertToDefaultSchema(ctx, table, []interface{}{entity})
		if err != nil {
			log.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	_, err := n.sb.Subscribe(subject, callback)

	return err
}

func (n *NatsWorker) subsNoBufferAsync(ctx context.Context, subject string, table string, wait bool) error {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.AsyncInsertToDefaultSchema(ctx, table, []interface{}{entity}, wait)
		if err != nil {
			log.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	_, err := n.sb.Subscribe(subject, callback)

	return err
}
