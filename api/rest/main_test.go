package rest_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go-dating-app/api/rest"
	"go-dating-app/app/repository"
	"go-dating-app/app/service"
	"go-dating-app/common/validation"
	"go-dating-app/config"
	"go-dating-app/database/migration"
	"go-dating-app/storage"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

//nolint:nolintlint,gochecknoglobals // testing globals.
var (
	testFlagsIntegration string
	db                   *sql.DB
	engine               *echo.Echo
	userRepo             repository.UserRepository
	authService          service.AuthService
	authHandler          *rest.AuthHandler
)

func TestMain(m *testing.M) {
	// Check if integration flag is set.
	testFlagsIntegration = os.Getenv("TEST_INTEGRATION") // valid values: "devcontainer". Other values will be ignored.
	if testFlagsIntegration != "devcontainer" {
		os.Exit(0)
	}

	// Setup server.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := setupEnv(logger)
	setupValidation(logger)
	setupDB(logger, cfg)
	teardownMigration(db, logger, cfg)
	setupMigration(db, logger)

	engine = echo.New()
	userRepo = repository.NewUserRepository(db)
	authService = service.NewAuthService(&userRepo)

	restApp := rest.App{
		Config: cfg,
		Logger: logger,
		Server: engine,
	}
	restServices := rest.Services{
		Auth: authService,
	}

	// Register handlers.
	authHandler = rest.NewAuthHandler(&restApp, &restServices)
	authHandler.Router()

	// Run tests.
	exitCode := m.Run()

	// Tear down: Drop all tables.
	teardownMigration(db, logger, cfg)

	// Exit.
	os.Exit(exitCode)
}

// cleanupDB: Deletes all rows from all tables.
func cleanupDB() {
	// Arrange them based on the order of relations in reverse.
	_, _ = db.Exec("DELETE FROM users;")
}

func doTestJSON(method, url string, dto any) (echo.Context, *httptest.ResponseRecorder) {
	reqBody, _ := json.Marshal(dto)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return engine.NewContext(req, rec), rec
}

func setupDB(logger *slog.Logger, cfg config.Config) {
	var err error
	db, err = storage.NewDB(cfg)
	if err != nil {
		logger.Error("Failed to load DB", "error", err)
		os.Exit(1)
	}
}

func setupEnv(logger *slog.Logger) config.Config {
	const dbDSN = "root:root@tcp(db:3306)/dating_test?charset=utf8mb4&parseTime=true"
	envPreset := map[string]string{
		"APP_ADDRESS": ":3000",
		"APP_NAME":    "go-dating-app-test",
		"APP_DEBUG":   "true",
		"DB_DSN":      dbDSN,
	}

	// Check for env vars or setup manually from the envPreset above.
	if os.Getenv("APP_NAME") == "" {
		for k, v := range envPreset {
			_ = os.Setenv(k, v)
		}
	}

	cfg, err := config.NewConfig(logger)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

func setupValidation(logger *slog.Logger) {
	if err := validation.NewValidation(); err != nil {
		logger.Error("Failed to load validation: ", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupMigration(db *sql.DB, logger *slog.Logger) {
	if err := migration.Up(db); err != nil {
		logger.Error("Failed to run migrations: ", slog.Any("error", err))
		os.Exit(1)
	}
}

func teardownMigration(db *sql.DB, logger *slog.Logger, cfg config.Config) {
	if err := migration.Reset(db, cfg); err != nil {
		logger.Error("Failed to reset migrations: ", slog.Any("error", err))
		os.Exit(1)
	}
}
