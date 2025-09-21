package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total des requêtes HTTP reçues",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Durée des requêtes HTTP",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	GRPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total des appels gRPC",
		},
		[]string{"method", "status"},
	)

	GRPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Durée des appels gRPC",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// RegisterCustomMetric permet à un service d’enregistrer ses propres métriques
func RegisterCustomMetric(c prometheus.Collector) {
	prometheus.MustRegister(c)
}
