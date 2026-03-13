package middleware

import (
	"net/http"
	"strconv"
	"time"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GrpcMetric = grpcprometheus.DefaultServerMetrics

	httpRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total_cust",
			Help: "HTPP Request Total",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_cust",
			Help:    "Duration of HTTP Request in Seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func MonitorMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			start := time.Now()

			wrapped := NewResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(wrapped.statusCode)

			httpRequestTotal.WithLabelValues(r.Method, path, status).Inc()
			httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		})
	}
}

func init() {
	GrpcMetric.EnableHandlingTimeHistogram()
}
