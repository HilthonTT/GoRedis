package handler

import (
	"goredis-server/internal/expiration"
	"time"
)

func (h *Handler) Expire(args []string) {
	key := args[1]

	seconds, _ := time.ParseDuration(args[2] + "s")
	expiration.SetExpiration(key, seconds)
}
