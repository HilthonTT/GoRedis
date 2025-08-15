package handler

import (
	"bufio"
	"fmt"
	"goredis-server/internal/data"
)

func (h *Handler) SRem(args []string) {
	if len(args) < 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SREM'\n"))
		return
	}

	key := args[1]
	members := args[2:]

	removedCount := h.DB.SRem(key, members...)

	w := bufio.NewWriter(h.conn)
	fmt.Fprintf(w, ":%d\r\n", removedCount)
	w.Flush()

	data.LogCommand("SREM", key, fmt.Sprintf("%v", members))
}
