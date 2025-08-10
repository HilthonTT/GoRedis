package handler

import "goredis-server/internal/messaging"

func (h *Handler) Publish(args []string) {
	topic := args[1]
	message := args[2]

	messaging.HandlePublish(topic, message)
	h.conn.Write([]byte("OK\n"))
}
