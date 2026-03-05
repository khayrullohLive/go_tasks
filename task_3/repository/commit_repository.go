package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type CommentRepository interface {
	FindByPostID(postID uint) ([]models.Comment, int64, error)
	FindByID(id uint) (*models.Comment, error)
	FindLike(userID, commentID uint) (*models.Like, error)

	// PostExists Post mavjudligini tekshirish
	PostExists(postID uint) bool

	// Create Youngish / yangilash / o'chirish
	Create(comment *models.Comment) error
	Update(comment *models.Comment, fields map[string]interface{}) error
	Delete(comment *models.Comment) error
	DeleteReplies(parentID uint) error

	// CreateLike Like
	CreateLike(like *models.Like) error
	DeleteLike(like *models.Like) error
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// FindByPostID Post kommentariylarini olish (faqat top-level, replies bilan)
func (r *commentRepository) FindByPostID(postID uint) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	r.db.Model(&models.Comment{}).
		Where("post_id = ? AND parent_id IS NULL", postID).
		Count(&total)

	result := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Author").
		Preload("Replies.Author").
		Order("created_at asc").
		Find(&comments)

	return comments, total, result.Error
}

// ID bo'yicha topish
func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	result := r.db.Preload("Author").First(&comment, id)
	return &comment, result.Error
}

// Post mavjudmi?
func (r *commentRepository) PostExists(postID uint) bool {
	var count int64
	r.db.Model(&models.Post{}).Where("id = ?", postID).Count(&count)
	return count > 0
}

// Yangi kommentariy saqlash
func (r *commentRepository) Create(comment *models.Comment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return err
	}
	// Yaratilgandan keyin Author ni preload qilib qaytarish
	return r.db.Preload("Author").First(comment, comment.ID).Error
}

// Kommentariyni yangilash
func (r *commentRepository) Update(comment *models.Comment, fields map[string]interface{}) error {
	if err := r.db.Model(comment).Updates(fields).Error; err != nil {
		return err
	}
	return r.db.Preload("Author").First(comment, comment.ID).Error
}

// Momentarily o'chirish (soft delete)
func (r *commentRepository) Delete(comment *models.Comment) error {
	return r.db.Delete(comment).Error
}

// Barchart replylarni o'chirish
func (r *commentRepository) DeleteReplies(parentID uint) error {
	return r.db.Where("parent_id = ?", parentID).Delete(&models.Comment{}).Error
}

// FindLike Like mavjudligini tekshirish
func (r *commentRepository) FindLike(userID, commentID uint) (*models.Like, error) {
	var like models.Like
	result := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like)
	return &like, result.Error
}

// Like qo'shish
func (r *commentRepository) CreateLike(like *models.Like) error {
	return r.db.Create(like).Error
}

// Like o'chirish
func (r *commentRepository) DeleteLike(like *models.Like) error {
	return r.db.Delete(like).Error
}
