package handler

import (
	"fmt"
	"goredis-server/internal/config"
)

func (h *Handler) Auth(args []string, cfg *config.Config) bool {
	if len(args) != 3 {
		fmt.Fprintln(h.conn, "ERR wrong number of arguments for 'AUTH'")
		return false
	}

	username := args[1]
	password := args[2]

	if username == cfg.Username && password == cfg.Password {
		fmt.Fprintln(h.conn, "OK")
		return true
	}

	fmt.Fprintln(h.conn, "ERR invalid username or password")
	return false
}
