package middleware

import (
	"net/http"
)

func WithTelemetry(next http.Handler) http.Handler {
	//return otelhttp.NewHandler(next, "dttmr-api")
	return next
}
