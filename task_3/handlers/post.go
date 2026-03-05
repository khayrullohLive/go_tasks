package handlers

import (
	"net/http"
	"strconv"
	"tasks/task_3/middleware"
	"tasks/task_3/models"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{db: db}
}

// GetPosts - Barcha postlarni olish (filtr, search, pagination bilan)
// GET /api/posts
func (h *PostHandler) GetPosts(c *gin.Context) {
	var query models.PostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Page size limitlash
	if query.PageSize > 50 {
		query.PageSize = 50
	}

	// Query builder
	db := h.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// Filtrlar
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	} else {
		db = db.Where("status = ?", models.PostPublished) // Default: faqat published
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

	// Jami soni
	var total int64
	db.Count(&total)

	// Saralash
	allowedSorts := map[string]bool{"created_at": true, "view_count": true, "title": true}
	sortBy := query.SortBy
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}
	sortOrder := query.SortOrder
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	// Ma'lumotlarni olish
	var posts []models.Post
	db.Order(sortBy + " " + sortOrder).
		Limit(query.PageSize).
		Offset(utils.GetOffset(query.Page, query.PageSize)).
		Find(&posts)

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       posts,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: utils.GetTotalPages(total, query.PageSize),
	})
}

// GetPost - Bitta postni olish
// GET /api/posts/:id
func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	result := h.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").Preload("Author").Preload("Replies.Author")
		}).
		First(&post, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	// Ko'rishlar sonini oshirish
	h.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	post.ViewCount++

	c.JSON(http.StatusOK, post)
}

// GetPostBySlug - Slug orqali postni olish
// GET /api/posts/slug/:slug
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var post models.Post
	result := h.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id IS NULL").Preload("Author").Preload("Replies.Author")
		}).
		Where("slug = ? AND status = ?", slug, models.PostPublished).
		First(&post)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	h.db.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	c.JSON(http.StatusOK, post)
}

// CreatePost - Yangi post yaratish
// POST /api/posts
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Noto'g'ri ma'lumotlar",
			Details: err.Error(),
		})
		return
	}

	// Excerpt avtomatik yaratish
	if req.Excerpt == "" {
		req.Excerpt = utils.GenerateExcerpt(req.Content, 200)
	}

	// Default status
	if req.Status == "" {
		req.Status = models.PostDraft
	}

	post := models.Post{
		Title:       req.Title,
		Slug:        utils.GenerateSlug(req.Title),
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		Status:      req.Status,
		AuthorID:    userID,
		CategoryID:  req.CategoryID,
		ReadingTime: utils.CalculateReadingTime(req.Content),
	}

	// Teglarni qo'shish (Many-to-Many)
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		h.db.Find(&tags, req.TagIDs)
		post.Tags = tags
	}

	if result := h.db.Create(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Post yaratishda xato"})
		return
	}

	// Preload qilib qaytarish
	h.db.Preload("Author").Preload("Category").Preload("Tags").First(&post, post.ID)

	c.JSON(http.StatusCreated, post)
}

// UpdatePost - Postni yangilash (faqat o'z posti yoki admin)
// PUT /api/posts/:id
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	isAdmin, _ := c.Get("is_admin")
	id := c.Param("id")

	var post models.Post
	if result := h.db.First(&post, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	// Ruxsat tekshirish
	if post.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu postni tahrirlash huquqingiz yo'q"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Yangilash
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
		updates["slug"] = utils.GenerateSlug(req.Title)
	}
	if req.Content != "" {
		updates["content"] = req.Content
		updates["reading_time"] = utils.CalculateReadingTime(req.Content)
	}
	if req.Excerpt != "" {
		updates["excerpt"] = req.Excerpt
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}

	h.db.Model(&post).Updates(updates)

	// Teglarni yangilash
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		h.db.Find(&tags, req.TagIDs)
		h.db.Model(&post).Association("Tags").Replace(tags)
	}

	h.db.Preload("Author").Preload("Category").Preload("Tags").First(&post, post.ID)
	c.JSON(http.StatusOK, post)
}

// DeletePost - Postni o'chirish (soft delete)
// DELETE /api/posts/:id
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	isAdmin, _ := c.Get("is_admin")
	id := c.Param("id")

	var post models.Post
	if result := h.db.First(&post, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	if post.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu postni o'chirish huquqingiz yo'q"})
		return
	}

	h.db.Delete(&post) // Soft delete - DeletedAt to'ldiriladi

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Post o'chirildi"})
}

// LikePost - Postga like bosish / bekor qilish (toggle)
// POST /api/posts/:id/like
func (h *PostHandler) LikePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Post mavjudligini tekshirish
	var post models.Post
	if result := h.db.First(&post, postID); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	// Avvalgi like ni tekshirish
	var like models.Like
	result := h.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like)

	if result.Error == nil {
		// Like mavjud -> o'chirish (toggle)
		h.db.Delete(&like)
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like bekor qilindi"})
	} else {
		// Like yo'q -> qo'shish
		pid := uint(postID)
		newLike := models.Like{UserID: userID, PostID: &pid}
		h.db.Create(&newLike)
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like qo'shildi"})
	}
}

// GetUserPosts - Foydalanuvchi postlarini olish
// GET /api/users/:id/posts
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	authorID := c.Param("id")

	var posts []models.Post
	var total int64

	h.db.Model(&models.Post{}).Where("author_id = ? AND status = ?", authorID, models.PostPublished).Count(&total)
	h.db.Where("author_id = ? AND status = ?", authorID, models.PostPublished).
		Preload("Category").Preload("Tags").
		Order("created_at desc").
		Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
	})
}
