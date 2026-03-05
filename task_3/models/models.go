package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password string `gorm:"not null" json:"-"` // json:"-" = JSON da ko'rinmaydi
	FullName string `gorm:"size:100" json:"full_name"`
	Bio      string `gorm:"type:text" json:"bio"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	// Relations
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:AuthorID" json:"comments,omitempty"`
	Likes    []Like    `json:"likes,omitempty"`

	// Follow system
	Followers []Follow `gorm:"foreignKey:FollowingID" json:"-"`
	Following []Follow `gorm:"foreignKey:FollowerID" json:"-"`
}

type Category struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null;size:100" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Color       string `gorm:"size:7;default:'#6366f1'" json:"color"` // Hex rang

	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

type Tag struct {
	BaseModel
	Name string `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Slug string `gorm:"uniqueIndex;not null;size:50" json:"slug"`

	// Many-to-Many: Post <-> Tag
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

type PostStatus string

const (
	PostDraft     PostStatus = "draft"
	PostPublished PostStatus = "published"
	PostArchived  PostStatus = "archived"
)

type Post struct {
	BaseModel
	Title       string     `gorm:"not null;size:255" json:"title"`
	Slug        string     `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"` // Qisqa tavsif
	CoverImage  string     `json:"cover_image"`
	Status      PostStatus `gorm:"type:varchar(20);default:'draft'" json:"status"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	ReadingTime int        `json:"reading_time"` // Daqiqada

	// Foreign Keys
	AuthorID   uint `gorm:"not null" json:"author_id"`
	CategoryID uint `json:"category_id"`

	// Relations (Preload qilganda keladi)
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category Category  `json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes,omitempty"`
}

// COMMENT MODEL (Nested comments - reply qo'llab-quvvatlaydi)
type Comment struct {
	BaseModel
	Content  string `gorm:"type:text;not null" json:"content"`
	IsEdited bool   `gorm:"default:false" json:"is_edited"`

	// Foreign Keys
	PostID   uint  `gorm:"not null" json:"post_id"`
	AuthorID uint  `gorm:"not null" json:"author_id"`
	ParentID *uint `json:"parent_id"` // nil = top-level comment

	// Relations
	Author  User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Post    Post      `gorm:"foreignKey:PostID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Likes   []Like    `gorm:"foreignKey:CommentID" json:"likes,omitempty"`
}

// LIKE MODEL - Post yoki Comment ga like bosish
type Like struct {
	BaseModel
	UserID    uint  `gorm:"not null" json:"user_id"`
	PostID    *uint `json:"post_id"`    // Post ga like
	CommentID *uint `json:"comment_id"` // Comment ga like

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// FOLLOW MODEL - Foydalanuvchilarni kuzatish
type Follow struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  uint      `gorm:"not null" json:"follower_id"`  // Kim kuzatmoqda
	FollowingID uint      `gorm:"not null" json:"following_id"` // Kimni kuzatmoqda
	CreatedAt   time.Time `json:"created_at"`

	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
