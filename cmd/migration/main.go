package main

import (
	"database/sql"
	"go-dating-app/config"
	"go-dating-app/database/migration"
	"go-dating-app/storage"
	"log/slog"
	"os"
)

func main() {
	var (
		err error
		db  *sql.DB
		cfg config.Config
	)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if cfg, err = config.NewConfig(logger); err != nil {
		logger.Error("Failed to load config: ", slog.Any("error", err))
		os.Exit(1)
	}
	if db, err = storage.NewDB(cfg); err != nil {
		logger.Error("Failed to load db: ", slog.Any("error", err))
		os.Exit(1)
	}

	if err = migration.Up(db); err != nil {
		logger.Error("Failed to run migrations: ", slog.Any("error", err))
		os.Exit(1)
	}
}
