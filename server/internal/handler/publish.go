package handler

import "goredis-server/internal/messaging"

func (h *Handler) Publish(args []string) {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'PUBLISH'\n"))
		return
	}

	topic := args[1]
	message := args[2]

	messaging.HandlePublish(topic, message)
	h.conn.Write([]byte("OK\n"))
}
