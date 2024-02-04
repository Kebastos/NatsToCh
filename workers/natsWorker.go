package workers

import (
	"context"
	client "github.com/Kebastos/NatsToCh/clients"
	"github.com/Kebastos/NatsToCh/config"
	"github.com/Kebastos/NatsToCh/log"
	"github.com/Kebastos/NatsToCh/models"
	"github.com/nats-io/nats.go"
	"time"
)

type ClickhouseStorage interface {
	BatchInsertToDefaultSchema(tableName string, data []interface{}, ctx context.Context) error
}

type NatsWorker struct {
	cfg []config.Subject
	sb  *client.NatsClient
	ch  ClickhouseStorage
}

func NewNatsWorker(cfg *config.Config, sb *client.NatsClient, ch ClickhouseStorage) *NatsWorker {
	return &NatsWorker{
		cfg: cfg.Subjects,
		sb:  sb,
		ch:  ch,
	}
}

func (n *NatsWorker) Start() error {
	for _, s := range n.cfg {
		if s.UseBuffer {
			return nil
		}

		err := n.subsNoBuffer(s.Name, s.TableName)
		if err != nil {
			return err
		}

		log.Infof("subscribed to %s\n with params: buffer=%t and table=%s", s.Name, s.UseBuffer, s.TableName)
	}

	return nil
}

func (n *NatsWorker) subsWithBuffer(subject string, table string) error {
	return nil
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
