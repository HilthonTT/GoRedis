package handler

import (
	"bufio"
	"fmt"
	"goredis-server/internal/data"
	"strings"
)

func (h *Handler) SRem(args []string) {
	if len(args) < 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SREM'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	members := args[2:]
	if len(members) == 0 {
		h.conn.Write([]byte("ERR no members provided for SREM\n"))
		return
	}

	removedCount := h.DB.SRem(key, members...)

	w := bufio.NewWriter(h.conn)
	fmt.Fprintf(w, ":%d\r\n", removedCount)
	w.Flush()

	data.LogCommand("SREM", key, fmt.Sprintf("%v", members))
}
