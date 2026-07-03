package router

import (
	"net/http"

	"github.com/robindittmar/dttmr-api/internal/api/handler"
	"github.com/robindittmar/dttmr-api/internal/api/middleware"
)

type Config struct{}

func NewMux(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.DefaultHandler)
	mux.HandleFunc("GET /health", handler.HealthHandler)

	var httpHandler http.Handler = mux
	httpHandler = middleware.WithTelemetry(httpHandler)

	return httpHandler
}
