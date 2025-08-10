package handler

import "goredis-server/internal/messaging"

func (h *Handler) Subscribe(args []string) {
	messaging.HandleSubscribe(h.conn, args[1])
	h.conn.Write([]byte("OK\n"))
}
