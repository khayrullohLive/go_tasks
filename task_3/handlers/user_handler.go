package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"tasks/task_3/middleware"
	"tasks/task_3/models"
	"tasks/task_3/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	repository repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repository: repo}
}

// GET /api/admin/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, total, err := h.repository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Foydalanuvchilarni olishda xato"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

// GET /api/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	user, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	postCount, followerCount, followingCount := h.repository.GetStats(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"user":            user,
		"post_count":      postCount,
		"follower_count":  followerCount,
		"following_count": followingCount,
	})
}

// POST /api/users/:id/follow  (toggle)
func (h *UserHandler) FollowUser(c *gin.Context) {
	followerID, _ := middleware.GetCurrentUserID(c)
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	// Maqsadli user mavjudligini tekshirish
	target, err := h.repository.FindByID(uint(targetID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	// O'zini kuzata olmaydi
	if followerID == target.ID {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "O'zingizni kuzata olmaysiz"})
		return
	}

	follow, err := h.repository.FindFollow(followerID, target.ID)
	if err == nil {
		// Mavjud → bekor qilish
		if err := h.repository.DeleteFollow(follow); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kuzatishdan chiqishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kuzatishdan chiqildi"})
	} else {
		// Yo'q → qo'shish
		if err := h.repository.CreateFollow(&models.Follow{
			FollowerID:  followerID,
			FollowingID: target.ID,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kuzatishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kuzatilmoqda"})
	}
}

// GET /api/users/:id/followers
func (h *UserHandler) GetFollowers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	followers, err := h.repository.FindFollowers(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kuzatuvchilarni olishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": followers, "total": len(followers)})
}

// GET /api/users/:id/following
func (h *UserHandler) GetFollowing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	following, err := h.repository.FindFollowing(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kuzatilayotganlarni olishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": following, "total": len(following)})
}
