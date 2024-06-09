package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"go-dating-app/config"
	"go-dating-app/database"

	"github.com/pressly/goose/v3"
)

var (
	ErrMigrationReset = errors.New("migrations can only be reset in debug mode")
)

func Up(db *sql.DB) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetBaseFS(database.EmbedMigrations)

	if err := goose.Up(db, "migration"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func Reset(db *sql.DB, cfg config.Config) error {
	if !cfg.App.Debug {
		return ErrMigrationReset
	}

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetBaseFS(database.EmbedMigrations)

	if err := goose.Reset(db, "migration"); err != nil {
		return fmt.Errorf("failed to reset db: %w", err)
	}
	return nil
}
