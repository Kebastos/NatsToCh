package workers

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
)

type ClickhouseWorker struct {
	cfg *config.Subject
	ch  ClickhouseStorage
	c   chan []interface{}
}

func NewClickhouseWorker(cfg *config.Subject, ch ClickhouseStorage, c chan []interface{}) *ClickhouseWorker {
	return &ClickhouseWorker{
		cfg: cfg,
		ch:  ch,
		c:   c,
	}
}

func (c *ClickhouseWorker) Start(ctx context.Context) {
	go func() {
		for {
			items := <-c.c
			if len(items) > 0 {
				err := c.ch.BatchInsertToDefaultSchema(ctx, c.cfg.TableName, items)
				if err != nil {
					log.Errorf("failed to insert data to clickhouse: %s", err)
				}
			}
		}
	}()
}
