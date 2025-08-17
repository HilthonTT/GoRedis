package main

import (
	"context"
	"errors"
	"goredis-shared/auth"
	"goredis-shared/config"
	"goredis-shared/env"
	"goredis-shared/redis"
	"goredis-web/internal/infrastructure/handler"
	"goredis-web/internal/infrastructure/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func main() {
	// Load environment variables
	httpAddr := env.GetEnv("HTTP_ADDR", ":8080")

	// Connect to Redis
	cli, err := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "guest",
		Password: "guest",
	})
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer cli.Close()

	// Init repositories
	todoItemRepository := repository.NewTodoItemRepository(cli)

	// Init auth store
	cookieStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})
	gothic.Store = cookieStore

	authService := auth.NewAuthService()

	// Setup Gin router
	router := gin.Default()

	h := handler.NewHandler(authService, todoItemRepository)
	router.GET("/login", h.HandleLogin)
	router.GET("/auth/:provider/callback", h.HandleCallbackFunction)
	router.GET("/auth/logout/:provider", h.HandleLogout)
	router.GET("/auth/:provider", h.HandleProviderLogin)
	router.GET("/", h.HandleHome)
	router.POST("/todos", h.HandleCreateTodo)
	router.GET("/todos/create", h.HandleCreateTodoPage)

	_, b, _, _ := runtime.Caller(0)                  // gets this file's path
	basePath := filepath.Join(filepath.Dir(b), "..") // go up one level
	staticPath := filepath.Join(basePath, "static")

	router.Static("/static", staticPath)

	// Create HTTP server
	srv := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	// Channel for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		log.Printf("🚀 Starting HTTP server on %s", httpAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Shutting down server...")

	// Gracefully shut down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	} else {
		log.Println("✅ Server exited gracefully")
	}
}
