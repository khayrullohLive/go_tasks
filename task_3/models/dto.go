package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type CreatePostRequest struct {
	Title      string     `json:"title" binding:"required,min=3,max=255"`
	Content    string     `json:"content" binding:"required"`
	Excerpt    string     `json:"excerpt"`
	CoverImage string     `json:"cover_image"`
	Status     PostStatus `json:"status"`
	CategoryID uint       `json:"category_id"`
	TagIDs     []uint     `json:"tag_ids"`
}

type UpdatePostRequest struct {
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Excerpt    string     `json:"excerpt"`
	CoverImage string     `json:"cover_image"`
	Status     PostStatus `json:"status"`
	CategoryID uint       `json:"category_id"`
	TagIDs     []uint     `json:"tag_ids"`
}

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"` // Reply uchun
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// Muvaffaqiyatli respose
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Xato response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type PostQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Search     string `form:"search"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	AuthorID   uint   `form:"author_id"`
	Status     string `form:"status"`
	SortBy     string `form:"sort_by,default=created_at"` // created_at, view_count, title
	SortOrder  string `form:"sort_order,default=desc"`    // asc, desc
}

type CommentQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}
