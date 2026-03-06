package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindByID(id uint) (*models.Category, error)
	Create(category *models.Category) error
	Delete(category *models.Category) error
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	result := r.db.Find(&categories)
	return categories, result.Error
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	result := r.db.Preload("Posts").First(&category, id)
	return &category, result.Error
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Delete(category *models.Category) error {
	return r.db.Delete(category).Error
}
