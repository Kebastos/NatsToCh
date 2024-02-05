package metrics

import (
	"github.com/Kebastos/NatsToCh/internal/config"
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

func MustRegister(config *config.Config) {
	namespace := config.Server.Namespace
	statusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "status_codes_total",
			Help:      "Distribution by status codes",
		},
		[]string{"user", "cluster", "cluster_user", "replica", "cluster_node", "code"},
	)
	concurrentQueries = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "concurrent_queries",
			Help:      "The number of concurrent queries at current time",
		},
		[]string{"user", "cluster", "cluster_user", "replica", "cluster_node"},
	)
	requestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "request_duration_seconds",
			Help:       "Request duration. Includes possible wait time in the queue",
			Objectives: map[float64]float64{0.5: 1e-1, 0.9: 1e-2, 0.99: 1e-3, 0.999: 1e-4, 1: 1e-5},
		},
		[]string{"user", "cluster", "cluster_user", "replica", "cluster_node"},
	)
	configSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "config_last_reload_successful",
		Help:      "Whether the last configuration reload attempt was successful.",
	})
	GotMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "got_message_count",
		Help:      "Total number of got messages",
	})
	InsertMessageCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "insert_message_count",
		Help:      "Total number of inserted messages",
	})

	prometheus.MustRegister(statusCodes, concurrentQueries, requestDuration, configSuccess, GotMessageCount, InsertMessageCount)
}
