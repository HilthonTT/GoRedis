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
		return nil, fmt.Errorf("failed to get todo IDs for user %s: %w", userID, err)
	}

	var todos []*domain.TodoItem
	for _, id := range ids {
		todo, err := r.GetByID(ctx, id)

		if err != nil {
			return nil, fmt.Errorf("failed to get todo %s: %w", id, err)
		}
		if todo != nil {
			todos = append(todos, todo)
		}
	}

	return todos, nil
}

func (r *todoItemRepository) GetByID(ctx context.Context, id string) (*domain.TodoItem, error) {
	key := fmt.Sprintf("todo:%s", id)

	data, err := r.client.Get(key)
	if err != nil {
		return nil, fmt.Errorf("redis GET %s failed: %w", key, err)
	}
	if data == "" {
		return nil, nil // not found
	}

	var todo domain.TodoItem
	if err := json.Unmarshal([]byte(data), &todo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal todo %s: %w", id, err)
	}

	return &todo, nil
}

func (r *todoItemRepository) Create(ctx context.Context, todo *domain.TodoItem) error {

	return r.save(ctx, todo)
}

func (r *todoItemRepository) Update(ctx context.Context, todo *domain.TodoItem) error {
	return r.save(ctx, todo)
}

func (r *todoItemRepository) save(ctx context.Context, todo *domain.TodoItem) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}

	key := fmt.Sprintf("todo:%s", todo.ID)
	bytes, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("failed to marshal todo: %w", err)
	}

	if err := r.client.Set(key, string(bytes)); err != nil {
		return fmt.Errorf("redis SET %s failed: %w", key, err)
	}

	// Track in user set
	userKey := fmt.Sprintf("user:%s:todos", todo.UserID)
	if err := r.client.SAdd(userKey, todo.ID); err != nil {
		return fmt.Errorf("failed to add todo %s to user set %s: %w", todo.ID, userKey, err)
	}

	return nil
}

func (r *todoItemRepository) Delete(ctx context.Context, todo *domain.TodoItem) error {
	key := fmt.Sprintf("todo:%s", todo.ID)
	if err := r.client.Delete(key); err != nil {
		return fmt.Errorf("failed to delete todo %s: %w", todo.ID, err)
	}

	userKey := fmt.Sprintf("user:%s:todos", todo.UserID)
	if err := r.client.SRem(userKey, todo.ID); err != nil {
		return fmt.Errorf("failed to remove todo %s from user set %s: %w", todo.ID, userKey, err)
	}

	return nil
}
