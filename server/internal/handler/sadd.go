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

	key := args[1]
	members := args[2:]

	added := h.DB.SAdd(key, members...)

	fmt.Fprintf(h.conn, "%d\n", added)

	data.LogCommand("SADD", key, strings.Join(members, " "))
}
