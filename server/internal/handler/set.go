package handler

import (
	"fmt"
	"goredis-server/internal/data"
)

func (h *Handler) Set(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(h.conn, "ERR wrong arguments")
		return
	}

	key, value := args[1], args[2]

	h.DB.Set(key, value)
	h.conn.Write([]byte("OK\n"))

	data.LogCommand("SET", key, value)
}
