package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordHashFailed   = errors.New("password hashing failed")
	ErrPasswordVerifyFailed = errors.New("password verify failed")
)

// Hash generate password hash from a string.
func Hash(s string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrPasswordHashFailed, err)
	}
	return string(hashedPassword), nil

}

// Verify password string with hashed password.
func Verify(plain, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false, fmt.Errorf("%w :%w", ErrPasswordVerifyFailed, err)
	}
	return true, nil
}
