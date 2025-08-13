package config

import (
	"goredis-shared/env"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost              string
	Port                    string
	CookiesAuthSecret       string
	CookiesAuthAgeInSeconds int
	CookiesAuthIsSecure     bool
	CookiesAuthIsHttpOnly   bool
	DiscordClientID         string
	DiscordClientSecret     string
	GithubClientID          string
	GithubClientSecret      string
}

const (
	TwoDaysInSeconds = 60 * 60 * 24 * 2
)

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		PublicHost:              env.GetEnv("PUBLIC_HOST", "http://localhost"),
		Port:                    env.GetEnv("PORT", "8080"),
		CookiesAuthSecret:       env.GetEnv("COOKIES_AUTH_SECRET", "some-very-secret-key"),
		CookiesAuthAgeInSeconds: env.GetEnvAsInt("COOKIES_AUTH_AGE_IN_SECONDS", TwoDaysInSeconds),
		CookiesAuthIsSecure:     env.GetEnvAsBool("COOKIES_AUTH_IS_SECURE", false),
		CookiesAuthIsHttpOnly:   env.GetEnvAsBool("COOKIES_AUTH_IS_HTTP_ONLY", false),
		DiscordClientID:         env.GetEnvOrError("DISCORD_CLIENT_ID"),
		DiscordClientSecret:     env.GetEnvOrError("DISCORD_CLIENT_SECRET"),
		GithubClientID:          env.GetEnvOrError("GITHUB_CLIENT_ID"),
		GithubClientSecret:      env.GetEnvOrError("GITHUB_CLIENT_SECRET"),
	}
}
