package service_test

import (
	"context"
	"errors"
	"go-dating-app/app/dto"
	"go-dating-app/app/entity"
	"go-dating-app/app/service"
	"go-dating-app/common/password"
	"go-dating-app/common/validation"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := validation.NewValidation(); err != nil {
		log.Fatal("Failed to load validation: ", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestAuth_Registration(t *testing.T) {
	mockAuthRepo := &service.MockAuthRepository{}
	authService := service.NewAuthService(mockAuthRepo)
	ctx := context.Background()

	type args struct {
		Email    string
		Password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "All valid", args: args{Email: "valid@example.com", Password: "validpassword"}, wantErr: false},
		{name: "Invalid email format", args: args{Email: "invalidexample.com", Password: "validpassword"}, wantErr: true},
		{name: "Invalid password length", args: args{Email: "invalid@example.com", Password: "inv"}, wantErr: true},
		{name: "All invalid", args: args{Email: "invalidexample", Password: "inv"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := dto.AuthRegistrationRequest{
				Email:    tt.args.Email,
				Password: tt.args.Password,
			}
			user := entity.User{
				Email:    input.Email,
				Password: input.Password,
			}

			if tt.wantErr == false {
				mockAuthRepo.On("Create", ctx, &user).Return(nil).Once()
			} else {
				mockAuthRepo.On("Create", ctx, &user).Return(errors.New("errors")).Once()
			}

			_, err := authService.Registration(ctx, input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Registration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthLogin(t *testing.T) {
	mockAuthRepo := &service.MockAuthRepository{}
	authService := service.NewAuthService(mockAuthRepo)
	ctx := context.Background()

	type args struct {
		Email            string
		Password         string
		PasswordMismatch bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "All valid", args: args{Email: "valid@example.com", Password: "validpassword"}},
		{name: "Password mismatch", args: args{Email: "valid@example.com", Password: "validpassword", PasswordMismatch: true}, wantErr: true},
		{name: "Invalid email format", args: args{Email: "invalidexample.com", Password: "validpassword"}, wantErr: true},
		{name: "Invalid password length", args: args{Email: "invalid@example.com", Password: "inv"}, wantErr: true},
		{name: "All invalid", args: args{Email: "invalidexample", Password: "inv"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := dto.AuthLoginRequest{
				Email:    tt.args.Email,
				Password: tt.args.Password,
			}
			user := entity.User{
				Email:    input.Email,
				Password: input.Password,
			}

			if !tt.args.PasswordMismatch {
				hashedPassword, err := password.Hash(user.Password)
				if err != nil {
					t.Errorf("Login() password hash failed %v", err)
				}
				user.Password = hashedPassword
			}

			if tt.wantErr == false || tt.args.PasswordMismatch {
				mockAuthRepo.On("FindByEmail", ctx, input.Email).Return(user, nil).Once()
			} else {
				mockAuthRepo.On("FindByEmail", ctx, input.Email).Return(user, errors.New("errors")).Once()
			}

			_, err := authService.Login(ctx, input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
