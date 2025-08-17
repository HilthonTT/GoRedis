package handler

import (
	"crypto/subtle"
	"fmt"
	"goredis-server/internal/config"
	"strings"
)

func (h *Handler) Auth(args []string, cfg *config.Config) bool {
	if len(args) != 3 {
		h.conn.Write([]byte("ERR wrong number of arguments for 'AUTH'\n"))
		return false
	}

	username := strings.TrimSpace(args[1])
	password := strings.TrimSpace(args[2])

	if username == "" || password == "" {
		h.conn.Write([]byte("ERR username and password cannot be empty\n"))
		return false
	}

	userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Username)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Password)) == 1

	if userMatch && passMatch {
		fmt.Fprintln(h.conn, "OK")
		return true
	}

	fmt.Fprintln(h.conn, "ERR invalid username or password")
	return false
}
