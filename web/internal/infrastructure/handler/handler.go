package handler

import (
	"goredis-shared/auth"
	"goredis-shared/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth  auth.AuthService
	redis *redis.Client
}

func NewHandler(auth auth.AuthService, redis *redis.Client) *Handler {
	return &Handler{
		auth:  auth,
		redis: redis,
	}
}

func redirect(c *gin.Context, location string) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", location)
		c.Status(http.StatusOK)
		return
	}
	c.Redirect(http.StatusFound, location)
}

func respondWithError(ctx *gin.Context, code int, message string) {
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.HTML(code, "error.tmpl", gin.H{
			"title":   "Error",
			"message": message,
		})
		return
	}

	problem := gin.H{
		"type":   "about:blank", // TODO: can use a URI for error type
		"title":  http.StatusText(code),
		"status": code,
		"detail": message,
	}

	ctx.JSON(code, problem)
}

func respondInternalError(ctx *gin.Context, message string) {
	respondWithError(ctx, http.StatusInternalServerError, message)
}
