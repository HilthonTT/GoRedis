package main

import (
	"goredis-shared/auth"
	"goredis-shared/config"
	"goredis-shared/redis"
	"goredis-web/internal/infrastructure/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cli, err := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "guest",
		Password: "guest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})
	authService := auth.NewAuthService()

	_ = handler.NewHandler(authService, cli)

	_ = gin.Default()

}
