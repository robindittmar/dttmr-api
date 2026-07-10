package domain

import (
	"context"
	"errors"
	"time"
)

type List struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ListRepository interface {
	CreateList(ctx context.Context, name string, userIDs []string) (*List, error)
}

type ListService struct {
	repo ListRepository
}

func NewListService(repo ListRepository) *ListService {
	return &ListService{repo: repo}
}

func (s *ListService) Create(ctx context.Context, name string, userIDs []string) (*List, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("users must have at least one associated user")
	}

	if name == "" {
		return nil, errors.New("list name must not be empty")
	}

	return s.repo.CreateList(ctx, name, userIDs)
}
