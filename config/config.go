package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	Name    string `env:"APP_NAME"`
	Address string `env:"APP_ADDRESS"`
	Debug   bool   `env:"APP_DEBUG"`
}

type dBConfig struct {
	DSN           string `env:"DB_DSN"`
	MaxOpenPool   int    `env:"DB_MAX_OPEN_POOL" envDefault:"10"`
	MaxIdlePool   int    `env:"DB_MAX_IDLE_POOL" envDefault:"10"`
	MaxIdleSecond int    `env:"DB_MAX_IDLE_SECOND" envDefault:"300"`
}

type Config struct {
	App appConfig
	DB  dBConfig
}

func NewConfig(logger *slog.Logger) (Config, error) {
	var cfg Config
	// Check if env vars are set.
	if os.Getenv("APP_ADDRESS") == "" {
		logger.Info("Env vars not found try loading .env file.")
		// Load .env file.
		if err := godotenv.Load(); err != nil {
			return cfg, fmt.Errorf("failed to load .env file: %w", err)
		}
		logger.Info(".env file successfully loaded.")
	}

	// Set all vars required unless they have default value.
	opts := env.Options{RequiredIfNoDef: true}

	// Parse env vars..
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return cfg, fmt.Errorf("failed to parse env vars: %w", err)
	}
	logger.Info("Env vars successfully parsed.")
	return cfg, nil
}
