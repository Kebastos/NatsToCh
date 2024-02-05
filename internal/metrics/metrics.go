package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	statusCodes        *prometheus.CounterVec
	concurrentQueries  *prometheus.GaugeVec
	requestDuration    *prometheus.SummaryVec
	configSuccess      prometheus.Gauge
	GotMessageCount    prometheus.Counter
	InsertMessageCount prometheus.Counter
)

func MustRegister() {
	GotMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "messages",
		Name:      "got_message_count",
		Help:      "Total number of got messages",
	})
	InsertMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "messages",
		Name:      "insert_message_count",
		Help:      "Total number of inserted messages",
	})

	prometheus.MustRegister(statusCodes, concurrentQueries, requestDuration, configSuccess, GotMessageCount, InsertMessageCount)
}
