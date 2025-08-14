package handler

import (
	"goredis-web/internal/presentation/views"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleLogin(c *gin.Context) {
	if err := views.Login().Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render login page: %v", err)
	}
}

func (h *Handler) HandleProviderLogin(c *gin.Context) {
	provider := c.Param("provider")
	c.Request = gothic.GetContextWithProvider(c.Request, provider)

	_, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err == nil {
		redirect(c, "/")
		return
	}

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *Handler) HandleCallbackFunction(c *gin.Context) {
	provider := c.Param("provider")
	c.Request = gothic.GetContextWithProvider(c.Request, provider)

	_, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		respondInternalError(c, "Authentication failed: "+err.Error())
		return
	}

	redirect(c, "/")
}

func (h *Handler) HandleLogout(c *gin.Context) {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		respondInternalError(c, "Authentication failed: "+err.Error())
		return
	}

	redirect(c, "/")
}
