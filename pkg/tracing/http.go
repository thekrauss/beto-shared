package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ajoute un middleware qui trace chaque requête HTTP
func HTTPTracingMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "http-request")
}
