package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tasks/task_3/middleware"
	"tasks/task_3/models"
	"tasks/task_3/repository"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	repository repository.PostRepository
	// ✅ db yo'q
}

func NewPostHandler(repo repository.PostRepository) *PostHandler {
	return &PostHandler{repository: repo}
}

// GET /api/posts
func (h *PostHandler) GetPosts(c *gin.Context) {
	var query models.PostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	posts, total, err := h.repository.FindAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Postlarni olishda xato"})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       posts,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: utils.GetTotalPages(total, query.PageSize),
	})
}

// GET /api/posts/:id
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	post, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	// Ko'rishlar sonini oshirish
	h.repository.IncrementViewCount(post.ID)
	post.ViewCount++

	c.JSON(http.StatusOK, post)
}

// GET /api/posts/slug/:slug
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	post, err := h.repository.FindBySlug(c.Param("slug"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	h.repository.IncrementViewCount(post.ID)
	post.ViewCount++

	c.JSON(http.StatusOK, post)
}

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

	if req.Excerpt == "" {
		req.Excerpt = utils.GenerateExcerpt(req.Content, 200)
	}
	if req.Status == "" {
		req.Status = models.PostDraft
	}

	post := &models.Post{
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

	// Teglarni olish va postga biriktirish
	if len(req.TagIDs) > 0 {
		tags, err := h.repository.FindTagsByIDs(req.TagIDs)
		if err == nil {
			post.Tags = tags
		}
	}

	if err := h.repository.Create(post); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Post yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// PUT /api/posts/:id
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	isAdmin, _ := c.Get("is_admin")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	post, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	if post.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu postni tahrirlash huquqingiz yo'q"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Faqat yuborilgan maydonlarni yangilash
	fields := map[string]interface{}{}
	if req.Title != "" {
		fields["title"] = req.Title
		fields["slug"] = utils.GenerateSlug(req.Title)
	}
	if req.Content != "" {
		fields["content"] = req.Content
		fields["reading_time"] = utils.CalculateReadingTime(req.Content)
	}
	if req.Excerpt != "" {
		fields["excerpt"] = req.Excerpt
	}
	if req.CoverImage != "" {
		fields["cover_image"] = req.CoverImage
	}
	if req.Status != "" {
		fields["status"] = req.Status
	}
	if req.CategoryID > 0 {
		fields["category_id"] = req.CategoryID
	}

	if err := h.repository.Update(post, fields); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Yangilashda xato"})
		return
	}

	// Teglarni yangilash
	if len(req.TagIDs) > 0 {
		tags, err := h.repository.FindTagsByIDs(req.TagIDs)
		if err == nil {
			h.repository.UpdateTags(post, tags)
		}
	}

	// Yangilangan postni qaytarish
	updated, _ := h.repository.FindByID(post.ID)
	c.JSON(http.StatusOK, updated)
}

// DELETE /api/posts/:id
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	isAdmin, _ := c.Get("is_admin")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	post, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	if post.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu postni o'chirish huquqingiz yo'q"})
		return
	}

	if err := h.repository.Delete(post); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "O'chirishda xato"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Post o'chirildi"})
}

// POST /api/posts/:id/like
func (h *PostHandler) LikePost(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	// Post mavjudligini tekshirish
	_, err = h.repository.FindByID(uint(postID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	like, err := h.repository.FindLike(userID, uint(postID))
	if err == nil {
		// Like mavjud → o'chirish
		if err := h.repository.DeleteLike(like); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Like o'chirishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like bekor qilindi"})
	} else {
		// Like yo'q → qo'shish
		pid := uint(postID)
		if err := h.repository.CreateLike(&models.Like{UserID: userID, PostID: &pid}); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Like qo'shishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like qo'shildi"})
	}
}

// GET /api/users/:id/posts
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	authorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	posts, total, err := h.repository.FindByAuthorID(uint(authorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Postlarni olishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"total": total,
	})
}
