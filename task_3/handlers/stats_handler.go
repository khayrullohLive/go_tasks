package handlers

import (
	"net/http"
	"tasks/task_3/models"
	"tasks/task_3/repository"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	repository repository.StatsRepository
}

func NewStatsHandler(repo repository.StatsRepository) *StatsHandler {
	return &StatsHandler{repository: repo}
}

// GET /api/admin/stats
func (h *StatsHandler) GetStats(c *gin.Context) {
	totalUsers, totalPosts, totalComments, totalLikes := h.repository.GetCounts()

	topPosts, err := h.repository.GetTopPosts(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Statistikani olishda xato"})
		return
	}

	recentUsers, err := h.repository.GetRecentUsers(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Statistikani olishda xato"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_users":    totalUsers,
		"total_posts":    totalPosts,
		"total_comments": totalComments,
		"total_likes":    totalLikes,
		"top_posts":      topPosts,
		"recent_users":   recentUsers,
	})
}
