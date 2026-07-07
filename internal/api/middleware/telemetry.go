package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func WithTelemetry(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "dttmr-api")
}
