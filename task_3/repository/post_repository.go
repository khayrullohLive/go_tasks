package repository

import (
	"tasks/task_3/models"
	"tasks/task_3/utils"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type PostRepository interface {
	// FindAll O'qish
	FindAll(query models.PostQuery) ([]models.Post, int64, error)
	FindByID(id uint) (*models.Post, error)
	FindBySlug(slug string) (*models.Post, error)
	FindByAuthorID(authorID uint) ([]models.Post, int64, error)
	FindLike(userID, postID uint) (*models.Like, error)
	FindTagsByIDs(tagIDs []uint) ([]models.Tag, error)

	// Create Yozish / yangilash / o'chirish
	Create(post *models.Post) error
	Update(post *models.Post, fields map[string]interface{}) error
	UpdateTags(post *models.Post, tags []models.Tag) error
	IncrementViewCount(postID uint) error
	Delete(post *models.Post) error

	// CreateLike Like
	CreateLike(like *models.Like) error
	DeleteLike(like *models.Like) error
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Barcha postlar — filtr, search, sort, pagination bilan
func (r *postRepository) FindAll(query models.PostQuery) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Page size limitlash
	if query.PageSize <= 0 || query.PageSize > 50 {
		query.PageSize = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	db := r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// Status filtri
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	} else {
		db = db.Where("status = ?", models.PostPublished)
	}

	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	if query.AuthorID > 0 {
		db = db.Where("author_id = ?", query.AuthorID)
	}

	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("title ILIKE ? OR content ILIKE ?", search, search)
	}

	if query.TagID > 0 {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", query.TagID)
	}

	// Jami son
	db.Count(&total)

	// Saralash — SQL injection oldini olish uchun whitelist
	allowedSorts := map[string]bool{"created_at": true, "view_count": true, "title": true}
	sortBy := query.SortBy
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}
	sortOrder := "desc"
	if query.SortOrder == "asc" {
		sortOrder = "asc"
	}

	result := db.Order(sortBy + " " + sortOrder).
		Limit(query.PageSize).
		Offset(utils.GetOffset(query.Page, query.PageSize)).
		Find(&posts)

	return posts, total, result.Error
}

// ID bo'yicha — kommentariylar bilan
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	result := r.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").
				Preload("Author").
				Preload("Replies.Author")
		}).
		First(&post, id)

	return &post, result.Error
}

// Slug bo'yicha — faqat published
func (r *postRepository) FindBySlug(slug string) (*models.Post, error) {
	var post models.Post
	result := r.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").
				Preload("Author").
				Preload("Replies.Author")
		}).
		Where("slug = ? AND status = ?", slug, models.PostPublished).
		First(&post)

	return &post, result.Error
}

// Muallif postlari
func (r *postRepository) FindByAuthorID(authorID uint) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	r.db.Model(&models.Post{}).
		Where("author_id = ? AND status = ?", authorID, models.PostPublished).
		Count(&total)

	result := r.db.
		Where("author_id = ? AND status = ?", authorID, models.PostPublished).
		Preload("Category").
		Preload("Tags").
		Order("created_at desc").
		Find(&posts)

	return posts, total, result.Error
}

// Tag IDlar bo'yicha teglarni topish
func (r *postRepository) FindTagsByIDs(tagIDs []uint) ([]models.Tag, error) {
	var tags []models.Tag
	result := r.db.Find(&tags, tagIDs)
	return tags, result.Error
}

// Yangi post yaratish
func (r *postRepository) Create(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return err
	}
	// Author, Category, Tags bilan qaytarish
	return r.db.Preload("Author").Preload("Category").Preload("Tags").First(post, post.ID).Error
}

// Postni yangilash
func (r *postRepository) Update(post *models.Post, fields map[string]interface{}) error {
	return r.db.Model(post).Updates(fields).Error
}

// Teglarni almashtirish (Many-to-Many)
func (r *postRepository) UpdateTags(post *models.Post, tags []models.Tag) error {
	return r.db.Model(post).Association("Tags").Replace(tags)
}

// Ko'rishlar sonini oshirish — atomic operatsiya
func (r *postRepository) IncrementViewCount(postID uint) error {
	return r.db.Model(&models.Post{}).
		Where("id = ?", postID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// Postni o'chirish (soft delete)
func (r *postRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}

// Like mavjudligini tekshirish
func (r *postRepository) FindLike(userID, postID uint) (*models.Like, error) {
	var like models.Like
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like)
	return &like, result.Error
}

// Like qo'shish
func (r *postRepository) CreateLike(like *models.Like) error {
	return r.db.Create(like).Error
}

// Like o'chirish
func (r *postRepository) DeleteLike(like *models.Like) error {
	return r.db.Delete(like).Error
}
