package config

import "goredis-server/internal/env"

type Config struct {
	BindAddr string
	Port     int
	Username string
	Password string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: env.GetString("BIND_ADDR", "127.0.0.1"),
		Username: env.GetString("AUTH_USERNAME", "guest"),
		Password: env.GetString("AUTH_PASSWORD", "guest"),
		Port:     env.GetInt("PORT", 6379),
	}
}
