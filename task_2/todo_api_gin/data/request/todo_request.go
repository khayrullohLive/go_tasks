package request

// CreateTodoRequest - yangi todo yaratish uchun so'rov
type CreateTodoRequest struct {
	Title       string `json:"title"       binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

// UpdateTodoRequest - todo yangilash uchun so'rov
type UpdateTodoRequest struct {
	Title       *string `json:"title"       binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Completed   *bool   `json:"completed"`
}
