package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type SearchRepository interface {
	Search(q string) (posts []models.Post, users []models.User, tags []models.Tag, err error)
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type searchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) Search(q string) ([]models.Post, []models.User, []models.Tag, error) {
	like := "%" + q + "%"

	var posts []models.Post
	if err := r.db.
		Where("(title ILIKE ? OR content ILIKE ?) AND status = ?", like, like, models.PostPublished).
		Preload("Author").
		Preload("Category").
		Limit(10).
		Find(&posts).Error; err != nil {
		return nil, nil, nil, err
	}

	var users []models.User
	if err := r.db.
		Where("username ILIKE ? OR full_name ILIKE ?", like, like).
		Limit(5).
		Find(&users).Error; err != nil {
		return nil, nil, nil, err
	}

	var tags []models.Tag
	if err := r.db.
		Where("name ILIKE ?", like).
		Limit(10).
		Find(&tags).Error; err != nil {
		return nil, nil, nil, err
	}

	return posts, users, tags, nil
}
