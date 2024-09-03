package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/models"
	"strconv"
)

type MetricInstrumentation interface {
	InsertMessageCountAdd(name string, count int)
}

type Client struct {
	cfg     config.CHConfig
	logger  *log.Log
	conn    clickhouse.Conn
	metrics MetricInstrumentation
}

func NewClickhouseClient(cfg *config.CHConfig, logger *log.Log, metrics MetricInstrumentation) *Client {
	return &Client{cfg: *cfg, logger: logger, metrics: metrics}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.cfg.Host + ":" + strconv.Itoa(c.cfg.Port)},
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
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	version, err := c.conn.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get ClickHouse server version: %w", err)
	}

	c.logger.Infof("connected to ClickHouse at %s with version %s", c.cfg.Host, version)
	return nil
}

func (c *Client) ConnStatus() driver.Stats {
	return c.conn.Stats()
}

func (c *Client) Close() error {
	c.logger.Infof("closing ClickHouse connection to %s", c.cfg.Host)
	return c.conn.Close()
}

func (c *Client) BatchInsert(ctx context.Context, tableName string, data []interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	query := fmt.Sprintf("INSERT INTO %s", tableName)

	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare batch for table %s: %w", tableName, err)
	}

	for _, row := range data {
		err = batch.AppendStruct(row)
		if err != nil {
			return fmt.Errorf("failed to append row to batch for table %s: %w", tableName, err)
		}
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("failed to send batch for table %s: %w", tableName, err)
	}

	c.metrics.InsertMessageCountAdd(tableName, len(data))

	return nil
}

func (c *Client) InsertAsync(ctx context.Context, tableName string, data *models.DefaultEntity) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ('%s', '%s', '%s', now(), '%s')", tableName, data.Id, data.ClientId, data.Subject, data.Content)

	if err := c.conn.AsyncInsert(ctx, query, false); err != nil {
		return fmt.Errorf("failed to prepare async insert for table %s: %w", tableName, err)
	}

	return nil
}
