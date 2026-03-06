package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type TagRepository interface {
	FindAll() ([]models.Tag, error)
	FindByID(id uint) (*models.Tag, error)
	Create(tag *models.Tag) error
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]models.Tag, error) {
	var tags []models.Tag
	result := r.db.Find(&tags)
	return tags, result.Error
}

// Posts.Author va Posts.Category bilan birga
func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	result := r.db.Preload("Posts.Author").Preload("Posts.Category").First(&tag, id)
	return &tag, result.Error
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}
