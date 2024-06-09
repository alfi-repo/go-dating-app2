package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-dating-app/app/entity"

	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	if err := user.OnSave(); err != nil {
		return err
	}

	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)",
		user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		// Check for duplicate entry.
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.ErrUserAlreadyExists
		}
		return fmt.Errorf("%w: %w", entity.ErrUserFailedToSave, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("%w: %w", entity.ErrUserFailedToGetID, err)
	}

	// Set user ID from database.
	user.ID = int(id)
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User

	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, password, created_at, updated_at, suspended_at FROM users WHERE email = ?",
		email,
	)
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.SuspendedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, entity.ErrUserNotFound
		}
		return user, fmt.Errorf("%w: %w", entity.ErrUserFailedToFind, err)
	}
	return user, nil
}
