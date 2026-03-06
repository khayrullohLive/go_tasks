package repository

import (
	"tasks/task_3/models"

	"gorm.io/gorm"
)

// ============================================================
// INTERFACE
// ============================================================

type UserRepository interface {
	FindAll() ([]models.User, int64, error)
	FindByID(id uint) (*models.User, error)
	GetStats(id uint) (postCount, followerCount, followingCount int64)

	FindFollow(followerID, followingID uint) (*models.Follow, error)
	FindFollowers(userID uint) ([]models.User, error)
	FindFollowing(userID uint) ([]models.User, error)
	CreateFollow(follow *models.Follow) error
	DeleteFollow(follow *models.Follow) error
}

// ============================================================
// IMPLEMENTATION
// ============================================================

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.User, int64, error) {
	var users []models.User
	var total int64
	r.db.Model(&models.User{}).Count(&total)
	result := r.db.Find(&users)
	return users, total, result.Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	return &user, result.Error
}

// Foydalanuvchi statistikasi — 3 ta count bir joyda
func (r *userRepository) GetStats(id uint) (postCount, followerCount, followingCount int64) {
	r.db.Model(&models.Post{}).Where("author_id = ? AND status = ?", id, models.PostPublished).Count(&postCount)
	r.db.Model(&models.Follow{}).Where("following_id = ?", id).Count(&followerCount)
	r.db.Model(&models.Follow{}).Where("follower_id = ?", id).Count(&followingCount)
	return
}

func (r *userRepository) FindFollow(followerID, followingID uint) (*models.Follow, error) {
	var follow models.Follow
	result := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&follow)
	return &follow, result.Error
}

func (r *userRepository) FindFollowers(userID uint) ([]models.User, error) {
	var follows []models.Follow
	r.db.Where("following_id = ?", userID).Preload("Follower").Find(&follows)

	users := make([]models.User, len(follows))
	for i, f := range follows {
		users[i] = f.Follower
	}
	return users, nil
}

func (r *userRepository) FindFollowing(userID uint) ([]models.User, error) {
	var follows []models.Follow
	r.db.Where("follower_id = ?", userID).Preload("Following").Find(&follows)

	users := make([]models.User, len(follows))
	for i, f := range follows {
		users[i] = f.Following
	}
	return users, nil
}

func (r *userRepository) CreateFollow(follow *models.Follow) error {
	return r.db.Create(follow).Error
}

func (r *userRepository) DeleteFollow(follow *models.Follow) error {
	return r.db.Delete(follow).Error
}
