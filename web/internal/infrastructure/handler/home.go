package handler

import (
	"goredis-web/internal/presentation/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleHome(c *gin.Context) {
	if err := views.Home().Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render home page: %v", err)
	}
}
