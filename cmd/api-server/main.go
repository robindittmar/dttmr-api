package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/robindittmar/dttmr-api/internal/api/router"
	"github.com/robindittmar/dttmr-api/internal/telemetry"
)

type Config struct {
	Environment  string
	Port         int
	OTLPEndpoint string
}

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
	setupLogging()

	slog.Info("Starting service", slog.String("service", serviceName), slog.String("version", serviceVersion))

	cfg := loadConfig()

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
			slog.Error("Failed to shutdown telemetry", err)
		}
	}()

	srv := makeServer(cfg)
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

func loadConfig() *Config {
	envFlag := flag.String("env", "development", "environment to use")
	portFlag := flag.Int("port", 8080, "port to listen on")
	otlpEndpointFlag := flag.String("otlp-endpoint", "localhost:4317", "otlp endpoint")

	flag.Parse()

	cfg := &Config{
		Environment:  *envFlag,
		Port:         *portFlag,
		OTLPEndpoint: *otlpEndpointFlag,
	}

	assignStringFromEnv("DTTMR_ENVIRONMENT", &cfg.Environment)
	assignIntFromEnv("DTTMR_PORT", &cfg.Port)
	assignStringFromEnv("DTTMR_OTLP_ENDPOINT", &cfg.OTLPEndpoint)

	return cfg
}

func assignStringFromEnv(key string, target *string) {
	if val, exists := os.LookupEnv(key); exists {
		*target = val
	}
}

func assignIntFromEnv(key string, target *int) {
	if val, exists := os.LookupEnv(key); exists {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			slog.Error("Failed to parse environment variable", slog.String("var", key), slog.Any("error", err))
		} else {
			*target = parsed
		}
	}
}

func makeServer(cfg *Config) *http.Server {
	routerConfig := router.Config{}
	mux := router.NewMux(routerConfig)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv
}
