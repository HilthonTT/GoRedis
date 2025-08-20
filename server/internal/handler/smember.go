package handler

import (
	"bufio"
	"fmt"
	"goredis-server/internal/data"
	"strings"
)

func (h *Handler) SMembers(args []string) {
	if len(args) != 2 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'SMEMBERS'\n"))
		return
	}

	key := strings.TrimSpace(args[1])
	if key == "" {
		h.conn.Write([]byte("ERR empty key is not allowed\n"))
		return
	}

	members := h.DB.SMembers(key)

	// Use a buffered writer to ensure all data is sent
	w := bufio.NewWriter(h.conn)

	if len(members) == 0 {
		fmt.Fprint(w, "*0\r\n")
		w.Flush()
		data.LogCommand("SMEMBERS", key, "")
		return
	}

	// RESP array header
	fmt.Fprintf(w, "*%d\r\n", len(members))

	// Each member as bulk string
	for _, member := range members {
		fmt.Fprintf(w, "$%d\r\n%s\r\n", len(member), member)
	}

	w.Flush() // Important! Push the data to the client

	data.LogCommand("SMEMBERS", key, "")
}
