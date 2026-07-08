package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/robindittmar/dttmr-api/internal/api/router"
	"github.com/robindittmar/dttmr-api/internal/telemetry"
)

func main() {
	serviceName := "dttmr-api"
	serviceVersion := "0.1.0"

	if err := run(serviceName, serviceVersion); err != nil {
		slog.Error("Service crashed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(serviceName string, serviceVersion string) error {
	_ = godotenv.Load(".env")

	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	traceHandler := &telemetry.TraceHandler{Handler: baseHandler}

	logger := slog.New(traceHandler)
	slog.SetDefault(logger)

	slog.Info("Starting service", slog.String("service", serviceName), slog.String("version", serviceVersion))

	shutdownTelemetry, err := telemetry.Init(context.Background(), serviceName, serviceVersion)
	if err != nil {
		slog.Error("Failed to initialize telemetry", err)
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := shutdownTelemetry(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown telemetry", err)
		}
	}()

	cfg := router.Config{}
	mux := router.NewMux(cfg)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

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
