package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"goredis-shared/redis"
	"goredis-web/internal/domain"
)

type todoItemRepository struct {
	client *redis.Client
}

func NewTodoItemRepository(client *redis.Client) *todoItemRepository {
	return &todoItemRepository{client}
}

func (r *todoItemRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.TodoItem, error) {
	userKey := fmt.Sprintf("user:%s:todos", userID)

	ids, err := r.client.SMembers(userKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo IDs for user: %w", err)
	}

	var todos []*domain.TodoItem
	for _, id := range ids {
		todo, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get todo: %w", err)
		}
		if todo != nil {
			todos = append(todos, todo)
		}
	}

	return todos, nil
}

func (r *todoItemRepository) GetByID(ctx context.Context, id string) (*domain.TodoItem, error) {
	key := fmt.Sprintf("todo:%v", id)

	data, err := r.client.Get(key)
	if err != nil {
		return nil, fmt.Errorf("redis GET failed: %w", err)
	}

	if data == "" {
		return nil, nil // not found
	}

	var todo domain.TodoItem
	if err := json.Unmarshal([]byte(data), &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *todoItemRepository) Create(ctx context.Context, todo *domain.TodoItem) error {
	key := fmt.Sprintf("todo:%v", todo.ID)
	bytes, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("failed to marshal todo: %w", err)
	}

	if err := r.client.Set(key, string(bytes)); err != nil {
		return fmt.Errorf("redis SET failed: %w", err)
	}

	// Track in user set
	userKey := fmt.Sprintf("user:%s:todos", todo.UserID)
	if err := r.client.SAdd(userKey, todo.ID); err != nil {
		return fmt.Errorf("failed to add todo ID to user set: %w", err)
	}

	return nil
}

func (r *todoItemRepository) Update(ctx context.Context, todo *domain.TodoItem) error {
	// Since Redis doesn't differentiate between create/update for SET,
	// we just overwrite the value.
	return r.Create(ctx, todo)
}

func (r *todoItemRepository) Delete(ctx context.Context, todo *domain.TodoItem) error {
	key := fmt.Sprintf("todo:%v", todo.ID)
	if err := r.client.Delete(key); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	userKey := fmt.Sprintf("user:%s:todos", todo.UserID)
	if err := r.client.SRem(userKey, todo.ID); err != nil {
		return fmt.Errorf("failed to remove todo ID from user set: %w", err)
	}

	return nil
}
