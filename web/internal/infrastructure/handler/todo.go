package handler

import (
	"goredis-shared/validation"
	"goredis-web/internal/domain"
	"goredis-web/internal/presentation/views"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/google/uuid"
)

const (
	MAX_MEMORY = 32 << 20
)

func (h *Handler) HandleCreateTodoPage(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		redirect(c, "/")
		return
	}

	buffer, err := renderHTML(c, views.AddTodoView(user))
	if err != nil {
		redirect(c, "/")
		return
	}

	respondHTML(c, buffer)
}

type CreateTodoRequest struct {
	Description string   `form:"description"`
	DueDate     string   `form:"dueDate"`
	Labels      []string `form:"labels"`
	Priority    int      `form:"priority"`
}

func (h *Handler) HandleCreateTodo(c *gin.Context) {

	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {

		respondWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := c.Request.ParseMultipartForm(MAX_MEMORY); err != nil {
		respondWithError(c, http.StatusBadRequest, "failed to parse form")
		return
	}

	var data CreateTodoRequest
	decoder := form.NewDecoder()
	if err := decoder.Decode(&data, c.Request.PostForm); err != nil {
		respondWithError(c, http.StatusBadRequest, "failed to decode form")
		return
	}

	v := validation.NewValidator()
	if err := v.Struct(data); err != nil {
		component := views.AddTodoView(user)
		component.Render(c, c.Writer)
		return
	}

	var dueDate time.Time
	if data.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", data.DueDate)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid due date format")
			return
		}
		dueDate = parsed
	}

	todo := &domain.TodoItem{
		ID:          uuid.NewString(),
		UserID:      user.UserID,
		Description: data.Description,
		DueDate:     dueDate,
		Labels:      data.Labels,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		Priority:    domain.Priority(data.Priority),
	}

	if err := h.repository.Create(c.Request.Context(), todo); err != nil {

		respondInternalError(c, "failed to create todo")
		return
	}

	// Respond
	// if c.GetHeader("HX-Request") == "true" {
	// 	component := views.TodoItemView(*todo)
	// 	component.Render(c, c.Writer)
	// 	return
	// }

	c.JSON(http.StatusCreated, todo)
}
