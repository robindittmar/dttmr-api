package router

import (
	"database/sql"
	"net/http"

	"github.com/robindittmar/dttmr-api/internal/api/handler"
	"github.com/robindittmar/dttmr-api/internal/api/middleware"
	"github.com/robindittmar/dttmr-api/internal/domain"
	"github.com/robindittmar/dttmr-api/internal/repository"
)

type Config struct {
	Database *sql.DB
}

func NewMux(cfg Config) http.Handler {
	listRepo := repository.NewListRepo(cfg.Database)
	listService := domain.NewListService(listRepo)

	listHandler := handler.ListHandler{ListService: listService}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.DefaultHandler)
	mux.HandleFunc("GET /health", handler.HealthHandler)

	mux.HandleFunc("POST /lists", listHandler.CreateList)

	var httpHandler http.Handler = mux
	httpHandler = middleware.WithMaxBytes(1024 * 64)(httpHandler)
	httpHandler = middleware.WithTelemetry(httpHandler)

	return httpHandler
}
