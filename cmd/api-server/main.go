package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/robindittmar/dttmr-api/internal/api/router"
	"github.com/robindittmar/dttmr-api/internal/config"
	"github.com/robindittmar/dttmr-api/internal/database"
	"github.com/robindittmar/dttmr-api/internal/telemetry"
)

func main() {
	serviceName := "dttmr-api"
	serviceVersion := "0.1.0"

	if err := run(serviceName, serviceVersion); err != nil {
		slog.Error("Service crashed")
		os.Exit(1)
	}
}

func run(serviceName string, serviceVersion string) error {
	_ = godotenv.Load(".env")
	setupLogging()

	slog.Info("Starting service", slog.String("service", serviceName), slog.String("version", serviceVersion))

	cfg := config.Load()

	telCfg := telemetry.Config{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Endpoint:       cfg.OTLPEndpoint,
		Environment:    cfg.Environment,
	}
	shutdownTelemetry, err := telemetry.Init(context.Background(), telCfg)
	if err != nil {
		slog.Error("Failed to initialize telemetry", err)
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := shutdownTelemetry(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown telemetry", slog.Any("error", err))
		}
	}()

	db, err := database.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to initialize database", slog.Any("error", err))
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("Failed to close database connection", slog.Any("error", err))
		}
	}()

	srv := makeServer(db, cfg.Port)
	go func() {
		slog.Info("Starting http server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start http server", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("Shutting down server...", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.Any("error", err))
		return err
	}

	slog.Info("Service shutdown successful!")
	return nil
}

func setupLogging() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	traceHandler := &telemetry.TraceHandler{Handler: baseHandler}

	logger := slog.New(traceHandler)
	slog.SetDefault(logger)
}

func makeServer(db *sql.DB, port int) *http.Server {
	routerConfig := router.Config{
		Database: db,
	}
	mux := router.NewMux(routerConfig)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv
}
