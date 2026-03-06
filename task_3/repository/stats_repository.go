package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type StatsRepository interface {
	GetCounts() (users, posts, comments, likes int64)
	GetTopPosts(limit int) ([]models.Post, error)
	GetRecentUsers(limit int) ([]models.User, error)
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

// 4 ta count ni bir metodda qaytarish
func (r *statsRepository) GetCounts() (users, posts, comments, likes int64) {
	r.db.Model(&models.User{}).Count(&users)
	r.db.Model(&models.Post{}).Count(&posts)
	r.db.Model(&models.Comment{}).Count(&comments)
	r.db.Model(&models.Like{}).Count(&likes)
	return
}

func (r *statsRepository) GetTopPosts(limit int) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.
		Where("status = ?", models.PostPublished).
		Preload("Author").
		Order("view_count desc").
		Limit(limit).
		Find(&posts)
	return posts, result.Error
}

func (r *statsRepository) GetRecentUsers(limit int) ([]models.User, error) {
	var users []models.User
	result := r.db.Order("created_at desc").Limit(limit).Find(&users)
	return users, result.Error
}
