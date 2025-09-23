package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// adds middleware that traces each HTTP request
func HTTPTracingMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http-request")
}
