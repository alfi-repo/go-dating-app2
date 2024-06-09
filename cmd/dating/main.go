package main

import (
	"database/sql"
	"go-dating-app/api/rest"
	"go-dating-app/app/repository"
	"go-dating-app/app/service"
	"go-dating-app/common/validation"
	"go-dating-app/config"
	"go-dating-app/storage"
	"log/slog"
	"os"
)

func main() {
	// Bootstrap.
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
	if err = validation.NewValidation(); err != nil {
		logger.Error("Failed to load validation: ", slog.Any("error", err))
		os.Exit(1)
	}

	// App: Repository.
	userRepo := repository.NewUserRepository(db)

	// App: Services
	authService := service.NewAuthService(&userRepo)

	// API: Rest.
	restApp := rest.App{
		Config: cfg,
		Logger: logger,
	}
	restServices := rest.Services{
		Auth: authService,
	}
	rest.StartServer(&restApp, &restServices)
}
