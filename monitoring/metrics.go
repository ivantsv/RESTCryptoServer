package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	CryptoOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "crypto_operations_total",
			Help: "Total number of crypto operations",
		},
		[]string{"operation", "symbol", "status"},
	)

	ExternalAPILatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_latency_seconds",
			Help:    "External API call latency",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"api", "endpoint"},
	)

	DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"operation", "table"},
	)

	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size_bytes",
			Help: "Current cache size in bytes",
		},
		[]string{"cache_type"},
	)

	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)
)

func RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func RecordCryptoOperation(operation, symbol, status string) {
	CryptoOperations.WithLabelValues(operation, symbol, status).Inc()
}

func RecordExternalAPICall(api, endpoint string, duration time.Duration) {
	ExternalAPILatency.WithLabelValues(api, endpoint).Observe(duration.Seconds())
}

func RecordDatabaseError(operation, table string) {
	DatabaseErrors.WithLabelValues(operation, table).Inc()
}

func UpdateCacheMetrics(cacheType string, size int64, hit bool) {
	CacheSize.WithLabelValues(cacheType).Set(float64(size))
	if hit {
		CacheHits.WithLabelValues(cacheType).Inc()
	} else {
		CacheMisses.WithLabelValues(cacheType).Inc()
	}
}