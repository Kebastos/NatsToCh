package metrics_test

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"testing"
	"time"
)

var m *metrics.Metrics
var err error

var (
	c         = make(chan prometheus.Metric)
	subject   = "test_subject"
	tableName = "test_table"
	cfg       = &config.Config{
		Subjects: []config.Subject{
			{
				Name:      subject,
				TableName: tableName,
				UseBuffer: true,
				BufferConfig: config.BufferConfig{
					MaxSize: 10,
					MaxWait: 10 * time.Second,
				},
			},
		},
	}
)

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestMetrics_MustRegister(t *testing.T) {
	m, err = metrics.NewMetrics(cfg)
	if err != nil {
		t.Errorf("failed to create metrics: %s", err)
	}
	m.MustRegister()
}

func TestNewMetricsWithEmptySubjects(t *testing.T) {
	cfg := &config.Config{}
	_, err := metrics.NewMetrics(cfg)
	if err == nil {
		t.Errorf("NewMetrics should return error")
	}
}

func TestMetrics_GotMessageCountInc(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	go m.GotMessageCountMap[subject].Collect(c)
	for {
		select {
		case <-c:
			log.Printf("got message from channel")
			return
		case <-ctx.Done():
			t.Errorf("did'n get message from cache")
		default:
			m.GotMessageCountInc(subject)
		}
	}
}

func TestMetrics_InsertMessageCountAdd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go m.InsertMessageCountMap[tableName].Collect(c)
	select {
	case count := <-c:
		log.Printf("got message from channel %s", count.Desc().String())
		return
	case <-ctx.Done():
		t.Errorf("did'n get message from cache")
	default:
		m.InsertMessageCountAdd(tableName, 100)
	}
}

func TestMetrics_QueueMessageCountInc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go m.QueueMessageCountMap[subject].Collect(c)
	select {
	case count := <-c:
		log.Printf("got message from channel %s", count.Desc().String())
		return
	case <-ctx.Done():
		t.Errorf("did'n get message from cache")
	default:
		m.QueueMessageCountInc(subject)
	}
}

func TestMetrics_QueueMessageCountDec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go m.QueueMessageCountMap[subject].Collect(c)
	select {
	case count := <-c:
		log.Printf("got message from channel %s", count.Desc().String())
		return
	case <-ctx.Done():
		t.Errorf("did'n get message from cache")
	default:
		m.QueueMessageCountDrain(subject)
	}
}
