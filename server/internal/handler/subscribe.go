package handler

import "goredis-server/internal/messaging"

func (h *Handler) Subscribe(args []string) {
	if len(args) < 2 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SUBSCRIBE'\n"))
		return
	}

	messaging.HandleSubscribe(h.conn, args[1])
}
