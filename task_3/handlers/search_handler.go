package handlers

import (
	"net/http"
	"tasks/task_3/models"
	"tasks/task_3/repository"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	repository repository.SearchRepository
}

func NewSearchHandler(repo repository.SearchRepository) *SearchHandler {
	return &SearchHandler{repository: repo}
}

// GET /api/search?q=golang
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Qidiruv so'zi kiritilmagan"})
		return
	}

	posts, users, tags, err := h.repository.Search(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Qidirishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": q,
		"posts": posts,
		"users": users,
		"tags":  tags,
	})
}
