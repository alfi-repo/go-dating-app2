package rest_test

import (
	"context"
	"go-dating-app/app/dto"
	"go-dating-app/app/entity"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createUser(email, password string) (entity.User, error) {
	var emptyUser entity.User
	ctx := context.Background()
	user, err := entity.NewUser(email, password)
	if err != nil {
		return emptyUser, err
	}

	if err = userRepo.Create(ctx, &user); err != nil {
		return emptyUser, err
	}
	return user, nil
}

func TestAuth_Registration(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		contains       string
	}{
		{
			name:           "can register",
			args:           args{email: "valid_email@example.com", password: "valid_password"},
			wantStatusCode: http.StatusCreated,
			contains:       "Registration success",
		},
		{
			name:           "can not register email exists",
			args:           args{email: "valid_email@example.com", password: "valid_password"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "email is already registered",
		},
		{
			name:           "can not register invalid email",
			args:           args{email: "valid_email_example.com", password: "ErrUserInvalidEmail"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "please check your input",
		},
		{
			name:           "can not register password too short",
			args:           args{email: "valid_email_2@example.com", password: "123"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "please check your input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := dto.AuthRegistrationRequest{
				Email:    tt.args.email,
				Password: tt.args.password,
			}
			ctx, rec := doTestJSON(http.MethodPost, "/auth/register", reqBody)
			if assert.NoError(t, authHandler.Registration(ctx)) {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				if tt.contains != "" {
					assert.Contains(t, rec.Body.String(), tt.contains)
				}
			}
		})
	}
	t.Cleanup(cleanupDB)
}

func TestAuthHandler_Login(t *testing.T) {
	// Create a valid user.
	validEmail := "valid_email@example.com"
	validPassword := "valid_password"
	if _, err := createUser(validEmail, validPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		contains       string
	}{
		{
			name:           "can login",
			args:           args{email: validEmail, password: validPassword},
			wantStatusCode: http.StatusOK,
			contains:       "Login success",
		},
		{
			name:           "can not login user not found",
			args:           args{email: "invalid_email@example.com", password: "valid_password"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "account not found or password invalid",
		},
		{
			name:           "can not login password invalid",
			args:           args{email: "valid_email@example.com", password: "invalid_password"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "account not found or password invalid",
		},
		{
			name:           "can not login password too short",
			args:           args{email: "valid_email@example.com", password: "123"},
			wantStatusCode: http.StatusBadRequest,
			contains:       "please check your input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := dto.AuthLoginRequest{
				Email:    tt.args.email,
				Password: tt.args.password,
			}
			ctx, rec := doTestJSON(http.MethodPost, "/auth/login", reqBody)
			if assert.NoError(t, authHandler.Login(ctx)) {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				if tt.contains != "" {
					assert.Contains(t, rec.Body.String(), tt.contains)
				}
			}
		})
	}
	t.Cleanup(cleanupDB)
}
