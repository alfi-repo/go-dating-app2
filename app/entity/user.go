package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"go-dating-app/common/password"
	"go-dating-app/common/validation"
	"time"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserFailedToSave      = errors.New("failed to save user")
	ErrUserFailedToGetID     = errors.New("failed to get user ID")
	ErrUserFailedToFind      = errors.New("failed to find user")
	ErrUserInvalidEmail      = errors.New("invalid email address")
	ErrUserInvalidPassword   = errors.New("invalid password")
	ErrUserPasswordHash      = errors.New("failed to hash password")
	ErrUserPasswordIncorrect = errors.New("incorrect password")
)

type User struct {
	ID          int
	Email       string
	Password    string // TODO should use custom type to prevent forgetting hashing when changing password
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SuspendedAt sql.NullTime
}

// CheckPassword check password. Return true if password is correct.
func (u *User) CheckPassword(plainPassword string) (bool, error) {
	ok, err := password.Verify(plainPassword, u.Password)
	if err != nil || !ok {
		return false, ErrUserPasswordIncorrect
	}
	return true, nil
}

// OnSave update timestamp on save.
func (u *User) OnSave() error {
	timeNow := time.Now().UTC()

	// New data or update
	if u.ID == 0 {
		hashedPassword, err := password.Hash(u.Password)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUserPasswordHash, err)
		}

		u.Password = hashedPassword
		u.CreatedAt = timeNow
		u.UpdatedAt = timeNow
	} else {
		u.UpdatedAt = timeNow
	}

	return nil
}

// NewUser create new user.
func NewUser(email, password string) (User, error) {
	var user User

	// Validate data.
	if err := validation.ValidateVar(email, "email,max=100"); err != nil {
		return user, ErrUserInvalidEmail
	}

	if err := validation.ValidateVar(password, "min=6"); err != nil {
		return user, ErrUserInvalidPassword
	}

	return User{
		Email:       email,
		Password:    password,
		SuspendedAt: sql.NullTime{Valid: false},
	}, nil
}
