package rest

import (
	"context"
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/service"
	"go-dating-app/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Config config.Config
	Logger *slog.Logger
	Server *echo.Echo
}

type Services struct {
	Auth service.AuthService
}

//nolint:gochecknoglobals // common errors shorthand.
var (
	ErrBadRequest     = dto.ErrorResponse{Message: "bad request"}
	ErrInternalServer = dto.ErrorResponse{Message: "internal server error"}
	ErrRequestTimeout = dto.ErrorResponse{Message: "request timeout, please try again"}
)

func StartServer(app *App, services *Services) {
	var err error
	// Setup server.
	restServer := echo.New()
	restServer.Use(middleware.Logger())
	restServer.Use(middleware.Recover())
	app.Server = restServer

	//  Register routers
	authHandler := NewAuthHandler(app, services)
	authHandler.Router()

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err = restServer.Start(app.Config.App.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			restServer.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = restServer.Shutdown(ctx); err != nil {
		restServer.Logger.Fatal(err)
	}
}
