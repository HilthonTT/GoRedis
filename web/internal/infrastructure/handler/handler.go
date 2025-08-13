package handler

import (
	"goredis-shared/auth"
	"goredis-shared/redis"
)

type Handler struct {
	auth  auth.AuthService
	redis *redis.Client
}

func NewHandler(auth auth.AuthService, redis *redis.Client) *Handler {
	return &Handler{
		auth:  auth,
		redis: redis,
	}
}
