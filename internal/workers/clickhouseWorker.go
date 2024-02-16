package workers

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
)

type ClickhouseWorker struct {
	cfg    *config.Subject
	ch     ClickhouseStorage
	c      chan []interface{}
	logger *log.Log
}

func NewClickhouseWorker(cfg *config.Subject, ch ClickhouseStorage, c chan []interface{}, logger *log.Log) *ClickhouseWorker {
	return &ClickhouseWorker{
		cfg:    cfg,
		ch:     ch,
		c:      c,
		logger: logger,
	}
}

func (c *ClickhouseWorker) Start(ctx context.Context) {
	go func() {
		for {
			items := <-c.c
			if len(items) > 0 {
				err := c.ch.BatchInsertToDefaultSchema(ctx, c.cfg.TableName, items)
				if err != nil {
					c.logger.Errorf("failed to insert data to clickhouse: %s", err)
				}
			}
		}
	}()
}
