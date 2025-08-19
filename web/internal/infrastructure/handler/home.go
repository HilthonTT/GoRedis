package handler

import (
	"goredis-web/internal/presentation/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleHome(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		redirect(c, "/login")
		return
	}

	todos, err := h.repository.GetByUserID(c.Request.Context(), user.UserID)

	if err := views.Home(user, todos).Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render home page: %v", err)
	}
}
