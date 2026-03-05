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

type CommentHandler struct {
	repository repository.CommentRepository
}

func NewCommentHandler(repo repository.CommentRepository) *CommentHandler {
	return &CommentHandler{repository: repo}
}

// GET /api/posts/:id/comments
func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	comments, total, err := h.repository.FindByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kommentariylarni olishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": total,
	})
}

// POST /api/posts/:id/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri post ID"})
		return
	}

	// Post mavjudligini tekshirish
	if !h.repository.PostExists(uint(postID)) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post topilmadi"})
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Reply bo'lsa parent mavjudligini tekshirish
	if req.ParentID != nil {
		_, err := h.repository.FindByID(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Parent kommentariy topilmadi"})
			return
		}
	}

	comment := &models.Comment{
		Content:  req.Content,
		PostID:   uint(postID),
		AuthorID: userID,
		ParentID: req.ParentID,
	}

	if err := h.repository.Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Kommentariy yaratishda xato"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// PUT /api/comments/:id
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	comment, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kommentariy topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	// Faqat o'z kommentariyini tahrirlash
	if comment.AuthorID != userID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu kommentariyni tahrirlash huquqingiz yo'q"})
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.repository.Update(comment, map[string]interface{}{
		"content":   req.Content,
		"is_edited": true,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Yangilashda xato"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DELETE /api/comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	isAdmin, _ := c.Get("is_admin")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	comment, err := h.repository.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Kommentariy topilmadi"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Xato yuz berdi"})
		return
	}

	if comment.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Bu kommentariyni o'chirish huquqingiz yo'q"})
		return
	}

	// Avval replylarni, keyin o'zini o'chirish
	if err := h.repository.DeleteReplies(comment.ID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Replylarni o'chirishda xato"})
		return
	}

	if err := h.repository.Delete(comment); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "O'chirishda xato"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Kommentariy o'chirildi"})
}

// POST /api/comments/:id/like
func (h *CommentHandler) LikeComment(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Noto'g'ri ID"})
		return
	}

	like, err := h.repository.FindLike(userID, uint(commentID))

	if err == nil {
		// Like mavjud → o'chirish
		if err := h.repository.DeleteLike(like); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Like o'chirishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like bekor qilindi"})
	} else {
		// Like yo'q → qo'shish
		cid := uint(commentID)
		if err := h.repository.CreateLike(&models.Like{UserID: userID, CommentID: &cid}); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Like qo'shishda xato"})
			return
		}
		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Like qo'shildi"})
	}
}
