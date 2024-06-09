package service

import (
	"context"
	"fmt"
	"go-dating-app/app/dto"
	"go-dating-app/app/entity"
)

type AuthRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	authRepo AuthRepository
}

func NewAuthService(authRepo AuthRepository) AuthService {
	return AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) Registration(ctx context.Context, dto dto.AuthRegistrationRequest) (entity.User, error) {
	newUser, err := entity.NewUser(dto.Email, dto.Password)
	if err != nil {
		return newUser, err
	}

	return newUser, s.authRepo.Create(ctx, &newUser)
}

func (s *AuthService) Login(ctx context.Context, dto dto.AuthLoginRequest) (AuthToken, error) {
	var authToken AuthToken
	user, err := s.authRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return authToken, err
	}

	// Check for password.
	passwordCorrect, err := user.CheckPassword(dto.Password)
	if err != nil || !passwordCorrect {
		return authToken, err
	}

	// Generate tokens.
	// TODO Change to proper token. i.e. JWT.
	authToken.AccessToken = fmt.Sprintf("dummy-access-token-%d", user.ID)
	authToken.RefreshToken = fmt.Sprintf("dummy-refresh-token-%d", user.ID)
	return authToken, nil
}
