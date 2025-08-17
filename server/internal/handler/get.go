package handler

import (
	"fmt"
	"goredis-server/internal/expiration"
	"strings"
	"time"
)

func (h *Handler) Get(args []string) {
	if len(args) != 2 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'GET'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

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
