package clients

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"strconv"
)

type ClickhouseClient struct {
	cfg  config.CHConfig
	conn clickhouse.Conn
}

func NewClickhouseClient(cfg *config.CHConfig) *ClickhouseClient {
	return &ClickhouseClient{cfg: *cfg}
}

func (c *ClickhouseClient) Connect() error {
	var err error
	c.conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.cfg.Host + ":" + strconv.Itoa(int(c.cfg.Port))},
		Auth: clickhouse.Auth{
			Database: c.cfg.Database,
			Username: c.cfg.User,
			Password: c.cfg.Password,
		},
		ConnMaxLifetime: c.cfg.ConnMaxLifetime,
		MaxOpenConns:    c.cfg.MaxOpenConns,
		MaxIdleConns:    c.cfg.MaxIdleConns,
	})

	if err != nil {
		return err
	}

	version, err := c.conn.ServerVersion()
	if err != nil {
		return err
	}

	log.Infof("connected to clickhouse %s with %s", c.cfg.Host, version)
	return nil
}

func (c *ClickhouseClient) BatchInsertToDefaultSchema(tableName string, data []interface{}, ctx context.Context) error {
	query := fmt.Sprintf("INSERT INTO %s", tableName)

	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		return err
	}

	for _, row := range data {
		err = batch.AppendStruct(row)
		if err != nil {
			return err
		}
	}

	err = batch.Send()
	if err != nil {
		return err
	}

	metrics.InsertMessageCount.Add(float64(len(data)))

	return nil
}

func (c *ClickhouseClient) AsyncInsertToDefaultSchema(tableName string, data []interface{}, wait bool, ctx context.Context) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES (@Subject, @CreateDateTime, @Content)", tableName)
	err := c.conn.AsyncInsert(ctx, query, wait, data...)
	if err != nil {
		return err
	}

	metrics.InsertMessageCount.Add(float64(len(data)))
	return nil
}