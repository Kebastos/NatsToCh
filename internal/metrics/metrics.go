package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GotMessageCount    prometheus.Counter
	InsertMessageCount prometheus.Counter
)

func MustRegister() {
	GotMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "messages",
		Name:      "got_message_total",
		Help:      "Total number of got messages",
	})
	InsertMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "messages",
		Name:      "insert_message_total",
		Help:      "Total number of inserted messages",
	})

	prometheus.MustRegister(GotMessageCount, InsertMessageCount)
}
