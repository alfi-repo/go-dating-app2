package rest

import (
	"context"
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/entity"
	"go-dating-app/common/validation"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	app      *App
	services *Services
}

func NewAuthHandler(app *App, services *Services) *AuthHandler {
	return &AuthHandler{
		app:      app,
		services: services,
	}
}

func (h *AuthHandler) Router() {
	routeAuth := h.app.Server.Group("/auth")
	routeAuth.POST("/register", h.Registration)
	routeAuth.POST("/login", h.Login)
}

func (h *AuthHandler) Registration(c echo.Context) error {
	var req dto.AuthRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrBadRequest)
	}

	// Validate request.
	if err := validation.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "please check your input",
			Errors:  validation.FormatStructErrors(err),
		})
	}

	// Handle user registration.
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()
	_, err := h.services.Auth.Registration(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return c.JSON(http.StatusServiceUnavailable, ErrRequestTimeout)
		case errors.Is(err, entity.ErrUserAlreadyExists):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "email is already registered",
			})
		default:
			h.app.Logger.Error("auth.registration", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Registration success",
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.AuthLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrBadRequest)
	}

	// Validate request.
	if err := validation.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "please check your input",
			Errors:  validation.FormatStructErrors(err),
		})
	}

	// Handle user login.
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()
	tokens, err := h.services.Auth.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return c.JSON(http.StatusServiceUnavailable, ErrRequestTimeout)
		case errors.Is(err, entity.ErrUserNotFound):
			fallthrough
		case errors.Is(err, entity.ErrUserPasswordIncorrect):
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "account not found or password invalid",
			})
		default:
			h.app.Logger.Error("auth.login", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login success",
		Data: dto.AuthLoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}
