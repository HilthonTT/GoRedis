package config

import "goredis-server/internal/env"

type Config struct {
	Port     int
	Username string
	Password string
}

func NewConfig() *Config {
	return &Config{
		Username: env.GetString("AUTH_USERNAME", "guest"),
		Password: env.GetString("AUTH_USERNAME", "guest"),
		Port:     env.GetInt("PORT", 6379),
	}
}
