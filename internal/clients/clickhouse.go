package clients

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/models"
	"reflect"
	"strconv"
	"strings"
)

type ClickhouseClient struct {
	cfg    config.CHConfig
	logger *log.Log
	conn   clickhouse.Conn
}

func NewClickhouseClient(cfg *config.CHConfig, logger *log.Log) *ClickhouseClient {
	return &ClickhouseClient{cfg: *cfg, logger: logger}
}

func (c *ClickhouseClient) Connect() error {
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
		return err
	}

	version, err := c.conn.ServerVersion()
	if err != nil {
		return err
	}

	c.logger.Infof("connected to clickhouse %s with %s", c.cfg.Host, version)
	return nil
}

func (c *ClickhouseClient) ConnStatus() driver.Stats {
	return c.conn.Stats()
}

func (c *ClickhouseClient) Close() error {
	c.logger.Infof("closing clickhouse %s", c.cfg.Host)
	return c.conn.Close()
}

func (c *ClickhouseClient) BatchInsertToDefaultTable(ctx context.Context, tableName string, data []interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

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

func (c *ClickhouseClient) AsyncInsertToDefaultTable(ctx context.Context, tableName string, data []interface{}, wait bool) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	t := reflect.TypeOf(models.DefaultTable{})

	// Prepare the query
	var fields []string
	var placeholders []string
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
		placeholders = append(placeholders, "?")
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(fields, ", "), strings.Join(placeholders, ", "))

	var insertData []interface{}
	for _, d := range data {
		v := reflect.ValueOf(d)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Type() == reflect.TypeOf(models.DefaultTable{}) {
			for _, field := range fields {
				insertData = append(insertData, v.FieldByName(field).Interface())
			}
		}
	}

	err := c.conn.AsyncInsert(ctx, query, wait, insertData...)
	if err != nil {
		return err
	}

	metrics.InsertMessageCount.Add(float64(len(data)))
	return nil
}
