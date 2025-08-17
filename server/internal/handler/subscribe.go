package handler

import (
	"goredis-server/internal/messaging"
	"strings"
)

func (h *Handler) Subscribe(args []string) {
	if len(args) < 2 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SUBSCRIBE'\n"))
		return
	}

	topic := strings.TrimSpace(args[1])
	if topic == "" {
		h.conn.Write([]byte("ERR empty topic is not allowed\n"))
		return
	}

	messaging.HandleSubscribe(h.conn, args[1])
}
