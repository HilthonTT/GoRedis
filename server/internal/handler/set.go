package handler

import (
	"goredis-server/internal/data"
)

func (h *Handler) Set(args []string) {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SET'\n"))
		return
	}

	key, value := args[1], args[2]

	h.DB.Set(key, value)
	h.conn.Write([]byte("OK\n"))

	data.LogCommand("SET", key, value)
}
