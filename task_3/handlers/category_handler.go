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

type CategoryHandler struct {
	repository repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repository: repo}
}

// GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kategoriyalarni olishda xato"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GET /api/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	category, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kategoriya topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// POST /api/admin/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        utils.GenerateSlug(req.Name),
		Description: req.Description,
		Color:       req.Color,
	}

	if err := h.repository.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kategoriya yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// DELETE /api/admin/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	category, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kategoriya topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	if err := h.repository.Delete(category); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "O'chirishda xato"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kategoriya o'chirildi"})
}
