package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE — Handler bu orqali gaplashadi
// ============================================================
// Interface yozishning foydasi:
//   1. Handler DB dan mustaqil bo'ladi
//   2. Test yozishda mock qilish oson (FakeAuthRepository)
//   3. Keyinchalik PostgreSQL -> MongoDB o'tsa, faqat shu fayl o'zgaradi

type AuthRepository interface {
	// User topish
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByEmailOrUsername(email, username string) (*models.User, error)

	// User yaratish / yangilash
	Create(user *models.User) error
	UpdateFields(user *models.User, fields map[string]interface{}) error
	UpdatePassword(userID uint, hashedPassword string) error
}

// ============================================================
// IMPLEMENTATION — Haqiqiy DB bilan ishlaydi
// ============================================================

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository Constructor
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindByEmailOrUsername Email yoki username bo'yicha topish (register uchun — mavjudligini tekshirish)
func (r *authRepository) FindByEmailOrUsername(email, username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ? OR username = ?", email, username).First(&user)
	if result.Error != nil {
		return nil, result.Error // gorm.ErrRecordNotFound bo'lsa — mavjud emas
	}
	return &user, nil
}

// FindByEmail Email bo'yicha topish (login uchun)
func (r *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByID ID bo'yicha topish (GetMe, UpdateProfile uchun)
func (r *authRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Posts").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create Yangi user yaratish
func (r *authRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// UpdateFields Bir necha maydonni yangilash (UpdateProfile uchun)
func (r *authRepository) UpdateFields(user *models.User, fields map[string]interface{}) error {
	return r.db.Model(user).Updates(fields).Error
}

// UpdatePassword Parolni yangilash
func (r *authRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
}
