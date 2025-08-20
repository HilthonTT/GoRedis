package handler

import (
	"strings"
)

func (h *Handler) Del(args []string) {
	if len(args) != 2 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'DEL'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	h.DB.Delete(key)

	h.conn.Write([]byte("OK\n"))
}
