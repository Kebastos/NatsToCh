package nats2ch

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/cache"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/models"
	"github.com/nats-io/nats.go"
)

func (n *Nats2Ch) callbackWithBuffer(ctx context.Context, cfg config.Subject) func(m *nats.Msg) {
	c := make(chan []interface{}, 1)
	cc := cache.New(&cfg, n.logger, c, n.metrics)
	cc.StartCleaner()
	n.startInsert(ctx, c, cfg.TableName)

	callback := func(m *nats.Msg) {
		entity := models.NewDefaultEntity(n.cfg.NATSConfig.ClientName, m.Subject, string(m.Data))

		cc.Set(entity)
	}

	return callback
}

func (n *Nats2Ch) callbackNoBuffer(ctx context.Context, table string) func(m *nats.Msg) {
	callback := func(m *nats.Msg) {
		entity := models.NewDefaultEntity(n.cfg.NATSConfig.ClientName, m.Subject, string(m.Data))

		err := n.ch.InsertAsync(ctx, table, entity)
		if err != nil {
			n.logger.Errorf("failed to insert data to clickhouse. %s", err)
		}
	}

	return callback
}
