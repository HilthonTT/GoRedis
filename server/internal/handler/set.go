package handler

import (
	"goredis-server/internal/data"
	"strings"
)

func (h *Handler) Set(args []string) {
	if len(args) < 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SET'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	value := strings.Join(args[2:], " ")

	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	h.DB.Set(key, value)
	h.conn.Write([]byte("OK\n"))

	data.LogCommand("SET", key, value)
}
