package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response - standart API javobi
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse - xato javobi
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   err,
	})
}

func NotFound(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   err,
	})
}

func InternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   err,
	})
}
