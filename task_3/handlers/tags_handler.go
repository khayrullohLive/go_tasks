package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tasks/task_3/models"
	"tasks/task_3/repository"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagHandler struct {
	repository repository.TagRepository
}

func NewTagHandler(repo repository.TagRepository) *TagHandler {
	return &TagHandler{repository: repo}
}

// GET /api/tags
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.repository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Teglarni olishda xato"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// POST /api/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	tag := &models.Tag{
		Name: req.Name,
		Slug: utils.GenerateSlug(req.Name),
	}

	if err := h.repository.Create(tag); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Teg yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GET /api/tags/:id/posts
func (h *TagHandler) GetTagPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	tag, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Teg topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":   tag,
		"posts": tag.Posts,
		"total": len(tag.Posts),
	})
}
