package handler

import (
	"goredis-server/internal/expiration"
	"time"
)

func (h *Handler) Expire(args []string) {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'EXPIRE'\n"))
		return
	}

	key := args[1]

	seconds, _ := time.ParseDuration(args[2] + "s")
	expiration.SetExpiration(key, seconds)
}
