package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// records HTTP metrics
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrapper captures the status code
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()

		//updates metrics
		HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rr.statusCode)).Inc()
		HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
