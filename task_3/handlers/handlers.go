package handlers

import (
	"net/http"
	"tasks/task_3/middleware"
	"tasks/task_3/models"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============================================================
// USER HANDLER
// ============================================================

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUsers - Barcha foydalanuvchilar (admin uchun)
// GET /api/admin/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []models.User
	var total int64

	h.db.Model(&models.User{}).Count(&total)
	h.db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

// GetUser - Bitta foydalanuvchi profili
// GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if result := h.db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
		return
	}

	// Statistika
	var postCount, followerCount, followingCount int64
	h.db.Model(&models.Post{}).Where("author_id = ? AND status = ?", id, models.PostPublished).Count(&postCount)
	h.db.Model(&models.Follow{}).Where("following_id = ?", id).Count(&followerCount)
	h.db.Model(&models.Follow{}).Where("follower_id = ?", id).Count(&followingCount)

	c.JSON(http.StatusOK, gin.H{
		"user":            user,
		"post_count":      postCount,
		"follower_count":  followerCount,
		"following_count": followingCount,
	})
}

// FollowUser - Foydalanuvchini kuzatish
// POST /api/users/:id/follow
func (h *UserHandler) FollowUser(c *gin.Context) {
	followerID, _ := middleware.GetCurrentUserID(c)
	followingID := c.Param("id")

	// O'zini kuzata olmasligi
	var followingUser models.User
	h.db.First(&followingUser, followingID)

	if followerID == followingUser.ID {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "O'zingizni kuzata olmaysiz"})
		return
	}

	// Mavjudligini tekshirish
	var follow models.Follow
	result := h.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow)

	if result.Error == nil {
		// Allaqachon kuzatyapti -> bekor qilish
		h.db.Delete(&follow)
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kuzatishdan chiqildi"})
	} else {
		// Yangi follow
		h.db.Create(&models.Follow{
			FollowerID:  followerID,
			FollowingID: followingUser.ID,
		})
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kuzatilmoqda"})
	}
}

// GetFollowers - Kuzatuvchilar ro'yxati
// GET /api/users/:id/followers
func (h *UserHandler) GetFollowers(c *gin.Context) {
	id := c.Param("id")
	var follows []models.Follow
	h.db.Where("following_id = ?", id).Preload("Follower").Find(&follows)

	followers := make([]models.User, len(follows))
	for i, f := range follows {
		followers[i] = f.Follower
	}

	c.JSON(http.StatusOK, gin.H{"data": followers, "total": len(followers)})
}

// GetFollowing - Kuzatilayotganlar ro'yxati
// GET /api/users/:id/following
func (h *UserHandler) GetFollowing(c *gin.Context) {
	id := c.Param("id")
	var follows []models.Follow
	h.db.Where("follower_id = ?", id).Preload("Following").Find(&follows)

	following := make([]models.User, len(follows))
	for i, f := range follows {
		following[i] = f.Following
	}

	c.JSON(http.StatusOK, gin.H{"data": following, "total": len(following)})
}

// ============================================================
// CATEGORY HANDLER
// ============================================================

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories - Barcha kategoriyalar
// GET /api/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	h.db.Find(&categories)
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetCategory - Bitta kategoriya
// GET /api/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if result := h.db.Preload("Posts").First(&category, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kategoriya topilmadi"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// CreateCategory - Yangi kategoriya (admin)
// POST /api/admin/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        utils.GenerateSlug(req.Name),
		Description: req.Description,
		Color:       req.Color,
	}

	if result := h.db.Create(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kategoriya yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// DeleteCategory - Kategoriyani o'chirish (admin)
// DELETE /api/admin/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if result := h.db.First(&category, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kategoriya topilmadi"})
		return
	}
	h.db.Delete(&category)
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kategoriya o'chirildi"})
}

// ============================================================
// TAG HANDLER
// ============================================================

type TagHandler struct {
	db *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

// GetTags - Barcha teglar
// GET /api/tags
func (h *TagHandler) GetTags(c *gin.Context) {
	var tags []models.Tag
	h.db.Find(&tags)
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// CreateTag - Yangi teg yaratish
// POST /api/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	tag := models.Tag{
		Name: req.Name,
		Slug: utils.GenerateSlug(req.Name),
	}

	if result := h.db.Create(&tag); result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Teg yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTagPosts - Teg bo'yicha postlar
// GET /api/tags/:id/posts
func (h *TagHandler) GetTagPosts(c *gin.Context) {
	id := c.Param("id")

	var tag models.Tag
	if result := h.db.Preload("Posts.Author").Preload("Posts.Category").First(&tag, id); result.Error != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Teg topilmadi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":   tag,
		"posts": tag.Posts,
		"total": len(tag.Posts),
	})
}

// ============================================================
// SEARCH HANDLER
// ============================================================

type SearchHandler struct {
	db *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

// Search - Umumiy qidiruv
// GET /api/search?q=golang
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Qidiruv so'zi kiritilmagan"})
		return
	}

	search := "%" + q + "%"

	var posts []models.Post
	h.db.Where("(title ILIKE ? OR content ILIKE ?) AND status = ?", search, search, models.PostPublished).
		Preload("Author").Preload("Category").
		Limit(10).Find(&posts)

	var users []models.User
	h.db.Where("username ILIKE ? OR full_name ILIKE ?", search, search).
		Limit(5).Find(&users)

	var tags []models.Tag
	h.db.Where("name ILIKE ?", search).Limit(10).Find(&tags)

	c.JSON(http.StatusOK, gin.H{
		"query": q,
		"posts": posts,
		"users": users,
		"tags":  tags,
	})
}

// ============================================================
// STATS HANDLER (admin dashboard uchun)
// ============================================================

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// GetStats - Umumiy statistika
// GET /api/admin/stats
func (h *StatsHandler) GetStats(c *gin.Context) {
	var totalUsers, totalPosts, totalComments, totalLikes int64

	h.db.Model(&models.User{}).Count(&totalUsers)
	h.db.Model(&models.Post{}).Count(&totalPosts)
	h.db.Model(&models.Comment{}).Count(&totalComments)
	h.db.Model(&models.Like{}).Count(&totalLikes)

	// Eng ko'p ko'rilgan postlar
	var topPosts []models.Post
	h.db.Where("status = ?", models.PostPublished).
		Preload("Author").
		Order("view_count desc").
		Limit(5).
		Find(&topPosts)

	// So'nggi foydalanuvchilar
	var recentUsers []models.User
	h.db.Order("created_at desc").Limit(5).Find(&recentUsers)

	c.JSON(http.StatusOK, gin.H{
		"total_users":    totalUsers,
		"total_posts":    totalPosts,
		"total_comments": totalComments,
		"total_likes":    totalLikes,
		"top_posts":      topPosts,
		"recent_users":   recentUsers,
	})
}
