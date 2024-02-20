package metrics

import (
	"fmt"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type Instrumentation interface {
	QueueMessageCountInc(name string)
	QueueMessageCountDrain(name string)
}

type Metrics struct {
	GotMessageCountMap    map[string]prometheus.Counter
	InsertMessageCountMap map[string]prometheus.Counter
	QueueMessageCountMap  map[string]prometheus.Gauge
}

func NewMetrics(cfg *config.Config) (*Metrics, error) {
	gotMap := map[string]prometheus.Counter{}
	insertMap := map[string]prometheus.Counter{}
	queueMap := map[string]prometheus.Gauge{}

	if cfg.Subjects == nil {
		return nil, fmt.Errorf("empty subjects configuration in config file")
	}

	for _, s := range cfg.Subjects {
		gotMap[s.Name] = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "messages",
			Name:      "got",
			Help:      "total number of got messages",
			ConstLabels: prometheus.Labels{
				"name": s.Name,
			},
		})

		insertMap[s.TableName] = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "messages",
			Name:      "inserted",
			Help:      "total number of inserted messages",
			ConstLabels: prometheus.Labels{
				"table": s.TableName,
			},
		})

		if s.UseBuffer {
			queueMap[s.Name] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "messages",
				Name:      "queue",
				Help:      "total number of messages in queue",
				ConstLabels: prometheus.Labels{
					"queue": s.Name,
				},
			})
		}
	}
	return &Metrics{
		GotMessageCountMap:    gotMap,
		InsertMessageCountMap: insertMap,
		QueueMessageCountMap:  queueMap,
	}, nil
}

func (m *Metrics) MustRegister() {
	cols := make([]prometheus.Collector, 0, len(m.GotMessageCountMap)+len(m.InsertMessageCountMap)+len(m.QueueMessageCountMap))
	for _, c := range m.GotMessageCountMap {
		cols = append(cols, c)
	}
	for _, c := range m.InsertMessageCountMap {
		cols = append(cols, c)
	}
	for _, c := range m.QueueMessageCountMap {
		cols = append(cols, c)
	}

	prometheus.MustRegister(cols...)
}

func (m *Metrics) GotMessageCountInc(name string) {
	c, ok := m.GotMessageCountMap[name]
	if !ok {
		log.Printf("get messages metric with name %s not found", name)
	}

	c.Inc()
}

func (m *Metrics) InsertMessageCountAdd(name string, count int) {
	c, ok := m.InsertMessageCountMap[name]
	if !ok {
		log.Printf("inser messages metric with name %s not found", name)
	}

	c.Add(float64(count))
}

func (m *Metrics) QueueMessageCountInc(name string) {
	c, ok := m.QueueMessageCountMap[name]
	if !ok {
		log.Printf("queue messages metric with name %s not found", name)
	}

	c.Inc()
}

func (m *Metrics) QueueMessageCountDrain(name string) {
	c, ok := m.QueueMessageCountMap[name]
	if !ok {
		log.Printf("queue messages metric with name %s not found", name)
	}

	c.Set(0)
}
