package handler

import (
	"fmt"
	"goredis-server/internal/data"
	"strings"
)

func (h *Handler) SAdd(args []string) {
	if len(args) < 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SADD'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	var members []string
	for _, m := range args[2:] {
		m = strings.TrimSpace(m)
		if m != "" {
			members = append(members, m)
		}
	}

	if len(members) == 0 {
		h.conn.Write([]byte("ERR no valid members to add\n"))
		return
	}

	added := h.DB.SAdd(key, members...)

	fmt.Fprintf(h.conn, "%d\n", added)

	data.LogCommand("SADD", key, strings.Join(members, " "))
}
