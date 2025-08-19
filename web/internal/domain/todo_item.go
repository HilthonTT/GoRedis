package domain

import (
	"context"
	"time"
)

type Priority int

const (
	Normal = 0
	Low    = 1
	Medium = 2
	High   = 3
	Top    = 4
)

type TodoItem struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Description string     `json:"description"`
	DueDate     time.Time  `json:"dueDate"`
	Labels      []string   `json:"labels"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Priority    Priority   `json:"priority"`
}

type TodoItemRepository interface {
	GetByUserID(ctx context.Context, userID string) ([]*TodoItem, error)
	GetByID(ctx context.Context, id string) (*TodoItem, error)
	Create(ctx context.Context, todo *TodoItem) error
	Update(ctx context.Context, todo *TodoItem) error
	Delete(ctx context.Context, todo *TodoItem) error
}
