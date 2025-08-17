package handler

import (
	"goredis-server/internal/messaging"
	"strings"
)

func (h *Handler) Publish(args []string) {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'PUBLISH'\n"))
		return
	}

	topic := strings.TrimSpace(args[1])
	message := strings.TrimSpace(args[2])

	if topic == "" {
		h.conn.Write([]byte("ERR empty topic is not allowed\n"))
		return
	}

	if message == "" {
		h.conn.Write([]byte("ERR empty message is not allowed\n"))
		return
	}

	messaging.HandlePublish(topic, message)

	h.conn.Write([]byte("OK\n"))
}
