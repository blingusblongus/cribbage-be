package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics.WritePrometheus(w, true)
}

func withMetrics(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	counter := func(status int) *metrics.Counter {
		return metrics.GetOrCreateCounter(fmt.Sprintf(`http_requests_total{endpoint=%q, status="%d"}`, endpoint, status))
	}
	duration := metrics.GetOrCreateHistogram(
		fmt.Sprintf(`http_request_duration_seconds{endpoint=%q}`, endpoint),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}

		next(rec, r)

		counter(rec.status).Inc()
		duration.UpdateDuration(start)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
