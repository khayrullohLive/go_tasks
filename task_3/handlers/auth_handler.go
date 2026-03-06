package handlers

import (
	"errors"
	"net/http"
	"tasks/task_3/models"
	"tasks/task_3/repository"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	jwtSecret  string
	repository repository.AuthRepository
}

func NewAuthHandler(jwtSecret string, repo repository.AuthRepository) *AuthHandler {
	return &AuthHandler{jwtSecret: jwtSecret, repository: repo}
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Noto'g'ri ma'lumotlar",
			Details: err.Error(),
		})
		return
	}

	// Email/username mavjudligini tekshirish — DB logika repositoryda
	existing, _ := h.repository.FindByEmailOrUsername(req.Email, req.Username)
	if existing != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Bu email yoki username allaqachon mavjud",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Parolni hash qilishda xato"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		IsAdmin:  req.IsAdmin,
	}

	// DB ga yozish — repositoryga topshirildi
	if err := h.repository.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Foydalanuvchi yaratishda xato"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.IsAdmin, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{Token: token, User: *user})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Noto'g'ri ma'lumotlar",
			Details: err.Error(),
		})
		return
	}

	// DB query — repositoryda
	user, err := h.repository.FindByEmail(req.Email)
	if err != nil {
		// Record topilmasa ham "noto'g'ri" deymiz — xavfsizlik uchun
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Email yoki parol noto'g'ri"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Hisob bloklangan"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Email yoki parol noto'g'ri"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.IsAdmin, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token yaratishda xato"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: *user})
}

// GET /api/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := h.repository.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /api/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.repository.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
		return
	}

	// Faqat to'ldirilgan maydonlarni yangilash
	fields := map[string]interface{}{}
	if req.FullName != "" {
		fields["full_name"] = req.FullName
	}
	if req.Bio != "" {
		fields["bio"] = req.Bio
	}
	if req.Avatar != "" {
		fields["avatar"] = req.Avatar
	}

	if err := h.repository.UpdateFields(user, fields); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Yangilashda xato"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Profil yangilandi", Data: user})
}

// PUT /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.repository.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Foydalanuvchi topilmadi"})
		return
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Eski parol noto'g'ri"})
		return
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Parolni hash qilishda xato"})
		return
	}

	if err := h.repository.UpdatePassword(userID, newHash); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Parolni yangilashda xato"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Parol muvaffaqiyatli o'zgartirildi"})
}
