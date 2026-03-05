package middleware

import (
	"net/http"
	"strings"
	"tasks/task_3/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - JWT tokenni tekshiradi
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header dan token olish
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header topilmadi"})
			c.Abort()
			return
		}

		// "Bearer TOKEN" formatini tekshirish
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token formati noto'g'ri: Bearer TOKEN bo'lishi kerak"})
			c.Abort()
			return
		}

		// Tokenni tekshirish
		claims, err := utils.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token yaroqsiz yoki muddati o'tgan"})
			c.Abort()
			return
		}

		// Kontekstga foydalanuvchi ma'lumotlarini saqlash
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

// AdminMiddleware - Faqat admin uchun
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu amalni bajarish uchun admin huquqi kerak"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalAuth - Token bo'lsa tekshiradi, bo'lmasa davom etadi
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := utils.ValidateToken(parts[1], jwtSecret); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("is_admin", claims.IsAdmin)
			}
		}

		c.Next()
	}
}

// CORS Middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger Middleware - So'rovlarni log qilish
func LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

// RateLimiter - Oddiy rate limiting (production'da Redis bilan ishlating)
func RateLimiter() gin.HandlerFunc {
	// Bu yerda oddiy implementatsiya
	// Real loyihada golang.org/x/time/rate yoki redis ishlatiladi
	return func(c *gin.Context) {
		c.Next()
	}
}

// GetCurrentUserID - Kontekstdan user ID olish helper
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
