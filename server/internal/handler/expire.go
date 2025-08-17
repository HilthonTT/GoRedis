package handler

import (
	"fmt"
	"goredis-server/internal/expiration"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) Expire(args []string) {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'EXPIRE'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	secondsStr := strings.TrimSpace(args[2])
	secondsInt, err := strconv.Atoi(secondsStr)
	if err != nil || secondsInt < 0 {
		h.conn.Write([]byte("ERR value is not an integer or out of range\n"))
		return
	}

	expiration.SetExpiration(key, time.Duration(secondsInt)*time.Second)

	fmt.Fprintln(h.conn, 1)
}
