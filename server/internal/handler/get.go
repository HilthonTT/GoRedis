package handler

import (
	"fmt"
	"goredis-server/internal/expiration"
	"time"
)

func (h *Handler) Get(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(h.conn, "ERR wrong arguments")
		return
	}

	key := args[1]

	expiry, hasExpiry := expiration.Expirations[key]
	now := time.Now()
	if hasExpiry && now.After(expiry) {
		h.DB.Delete(key)
		expiration.RemoveExpiration(key)
		fmt.Fprintln(h.conn, "(nil)")
		return
	}

	val, ok := h.DB.Get(key)
	if !ok {
		fmt.Fprintln(h.conn, "(nil)")
	} else {
		fmt.Fprintln(h.conn, val)
	}
}
