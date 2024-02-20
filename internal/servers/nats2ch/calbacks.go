package nats2ch

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/cache"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/models"
	"github.com/nats-io/nats.go"
	"time"
)

func (n *Nats2Ch) callbackWithBuffer(ctx context.Context, cfg config.Subject) func(m *nats.Msg) {
	c := make(chan []interface{}, 1)
	cc := cache.New(&cfg, n.logger, c, n.metrics)
	cc.StartCleaner()
	n.startInsert(ctx, c, cfg.TableName)

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

func (n *Nats2Ch) callbackNoBuffer(ctx context.Context, table string) func(m *nats.Msg) {
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

func (n *Nats2Ch) callbackNoBufferAsync(ctx context.Context, table string, wait bool) func(m *nats.Msg) {
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
