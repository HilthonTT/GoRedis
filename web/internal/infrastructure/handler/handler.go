package handler

import (
	"bytes"
	"goredis-shared/auth"
	"goredis-web/internal/domain"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth       auth.AuthService
	repository domain.TodoItemRepository
}

func NewHandler(auth auth.AuthService, repository domain.TodoItemRepository) *Handler {
	return &Handler{
		auth:       auth,
		repository: repository,
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

func respondNotfound(ctx *gin.Context, message string) {
	respondWithError(ctx, http.StatusNotFound, message)
}

func respondHTML(c *gin.Context, buffer []byte) {
	c.Data(200, "text/html; charset=utf-8", buffer)
}

func renderHTML(c *gin.Context, template templ.Component) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	err := template.Render(c, buffer)
	if err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}
