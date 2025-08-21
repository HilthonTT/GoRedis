package handler

import (
	"goredis-shared/validation"
	"goredis-web/internal/domain"
	"goredis-web/internal/presentation/views"
	"log"
	"net/http"
	"strings"
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

func (h *Handler) HandleTodoView(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		redirect(c, "/login")
		return
	}

	todoId := c.Param("id")

	ctx := c.Request.Context()
	todo, err := h.repository.GetByID(ctx, todoId)
	if err != nil {
		respondInternalError(c, "failed to load todo")
		return
	}

	if todo == nil {
		redirect(c, "/")
		return
	}

	buffer, err := renderHTML(c, views.TodoView(user, todo))
	if err != nil {
		redirect(c, "/")
		return
	}

	respondHTML(c, buffer)
}

func (h *Handler) HandleDeleteTodo(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		redirect(c, "/login")
		return
	}

	todoId := c.Param("id")

	ctx := c.Request.Context()
	todo, err := h.repository.GetByID(ctx, todoId)
	if err != nil {
		respondInternalError(c, "failed to load todo")
		return
	}

	if todo == nil {
		respondNotfound(c, "todo not found")
		return
	}

	if todo.UserID != user.UserID {
		respondWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.repository.Delete(ctx, todo); err != nil {
		respondInternalError(c, "failed to delete todo")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) OpenEditTodoModal(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		redirect(c, "/login")
		return
	}

	todoId := c.Param("id")

	ctx := c.Request.Context()
	todo, err := h.repository.GetByID(ctx, todoId)
	if err != nil {
		respondInternalError(c, "failed to delete todo")
		return
	}

	if todo == nil {
		respondNotfound(c, "todo not found")
		return
	}

	buffer, err := renderHTML(c, views.EditTodoModal(user, todo))
	if err != nil {
		redirect(c, "/")
		return
	}

	respondHTML(c, buffer)
}

type UpdateTodoRequest struct {
	Description string   `form:"description" validate:"required"`
	DueDate     string   `form:"dueDate"`
	Labels      []string `form:"labels"`
	Priority    string   `form:"priority" validate:"oneof=low medium high"`
}

func (h *Handler) HandleUpdateTodo(c *gin.Context) {
	user, err := h.auth.GetSessionUser(c.Request)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := c.Request.ParseMultipartForm(MAX_MEMORY); err != nil {

		respondWithError(c, http.StatusBadRequest, "failed to parse form")
		return
	}

	var data UpdateTodoRequest
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

	todoId := c.Param("id")
	ctx := c.Request.Context()

	todo, err := h.repository.GetByID(ctx, todoId)
	if err != nil {
		respondNotfound(c, "todo not found")
		return
	}

	var priority domain.Priority
	switch strings.ToLower(data.Priority) {
	case "low":
		priority = domain.Low
	case "medium":
		priority = domain.Medium
	case "high":
		priority = domain.High
	default:
		log.Printf("Invalid priority: %s", data.Priority)
		respondWithError(c, http.StatusBadRequest, "invalid priority")
		return
	}

	todo.Priority = domain.Priority(priority)

	todo.Description = data.Description
	todo.DueDate = dueDate

	if len(data.Labels) == 1 {
		labels := strings.Split(strings.TrimSpace(data.Labels[0]), ",")
		todo.Labels = make([]string, 0, len(labels))
		for _, label := range labels {
			if trimmed := strings.TrimSpace(label); trimmed != "" {
				todo.Labels = append(todo.Labels, trimmed)
			}
		}
	} else {
		todo.Labels = data.Labels
	}

	if err := h.repository.Update(ctx, todo); err != nil {
		respondInternalError(c, "failed to update todo")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	views.TodoItem(todo).Render(c, c.Writer)
	views.CloseModal().Render(c, c.Writer)
}

func (h *Handler) CloseModal(c *gin.Context) {
	buffer, err := renderHTML(c, views.CloseModal())
	if err != nil {
		respondInternalError(c, "failed to render close modal")
		return
	}
	respondHTML(c, buffer)
}
