package handler

import (
	"goredis-server/internal/cache"
	"net"
)

type Handler struct {
	DB   cache.ShardMap
	conn net.Conn
}

func NewHandler() *Handler {
	return &Handler{
		DB:   cache.NewShardMap(16),
		conn: nil,
	}
}

func (h *Handler) SetConn(conn net.Conn) {
	h.conn = conn
}
