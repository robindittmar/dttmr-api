package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/robindittmar/dttmr-api/internal/domain"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) GetByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	user := &domain.AuthUser{}

	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, name, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
