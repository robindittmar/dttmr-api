package config

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Environment  string
	Port         int
	OTLPEndpoint string
	DatabaseURL  string
}

func Load() *Config {
	envFlag := flag.String("env", "development", "environment to use")
	portFlag := flag.Int("port", 8080, "port to listen on")
	otlpEndpointFlag := flag.String("otlp-endpoint", "localhost:4317", "otlp endpoint")
	databaseUrlFlag := flag.String("database-url", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "database connection string")

	flag.Parse()

	cfg := &Config{
		Environment:  *envFlag,
		Port:         *portFlag,
		OTLPEndpoint: *otlpEndpointFlag,
		DatabaseURL:  *databaseUrlFlag,
	}

	assignStringFromEnv("DTTMR_ENVIRONMENT", &cfg.Environment)
	assignIntFromEnv("DTTMR_PORT", &cfg.Port)
	assignStringFromEnv("DTTMR_OTLP_ENDPOINT", &cfg.OTLPEndpoint)
	assignStringFromEnv("DTTMR_DATABASE_URL", &cfg.DatabaseURL)

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
