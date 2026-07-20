package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/robindittmar/dttmr-api/internal/domain"
)

type ListRepo struct {
	db *sql.DB
}

func NewListRepo(db *sql.DB) *ListRepo {
	return &ListRepo{db: db}
}

func (r *ListRepo) CreateList(ctx context.Context, name string, userIDs []string) (*domain.List, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			slog.Error("failed to rollback transaction", slog.Any("error", err))
		}
	}()

	list := &domain.List{Name: name}

	err = tx.QueryRowContext(ctx,
		"INSERT INTO lists (name) VALUES ($1) RETURNING id, created_at",
		name,
	).Scan(&list.ID, &list.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert list: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO list_users (list_id, user_id) VALUES ($1, $2)")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare user/list association statement: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			slog.Error("failed to close user/list association statement", slog.Any("error", err))
		}
	}()

	for _, userID := range userIDs {
		if _, err = stmt.ExecContext(ctx, list.ID, userID); err != nil {
			return nil, fmt.Errorf("failed to insert user/list association: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return list, err
}
