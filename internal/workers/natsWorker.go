package workers

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/cache"
	client "github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/models"
	"github.com/nats-io/nats.go"
	"time"
)

type ClickhouseStorage interface {
	BatchInsertToDefaultSchema(tableName string, data []interface{}, ctx context.Context) error
	AsyncInsertToDefaultSchema(tableName string, data []interface{}, wait bool, ctx context.Context) error
}

type NatsWorker struct {
	cfg *config.Config
	sb  *client.NatsClient
	ch  ClickhouseStorage
}

func NewNatsWorker(cfg *config.Config, sb *client.NatsClient, ch ClickhouseStorage) *NatsWorker {
	return &NatsWorker{
		cfg: cfg,
		sb:  sb,
		ch:  ch,
	}
}

func (n *NatsWorker) Start(ctx context.Context) error {
	for _, s := range n.cfg.Subjects {
		var err error

		if s.UseBuffer {
			err = n.subsWithBuffer(ctx, s.Name, s)
		} else if s.Async {
			err = n.subsNoBufferAsync(s.Name, s.TableName, s.AsyncInsertConfig.Wait)
		} else {
			err = n.subsNoBuffer(s.Name, s.TableName)
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

func (n *NatsWorker) subsNoBuffer(subject string, table string) error {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.BatchInsertToDefaultSchema(table, []interface{}{entity}, context.Background())
		if err != nil {
			log.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	_, err := n.sb.Subscribe(subject, callback)

	return err
}

func (n *NatsWorker) subsNoBufferAsync(subject string, table string, wait bool) error {
	callback := func(m *nats.Msg) {
		msg := string(m.Data)

		entity := &models.DefaultTable{
			Subject:        m.Subject,
			CreateDateTime: time.Now(),
			Content:        msg,
		}

		err := n.ch.AsyncInsertToDefaultSchema(table, []interface{}{entity}, wait, context.Background())
		if err != nil {
			log.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	_, err := n.sb.Subscribe(subject, callback)

	return err
}
