package main

import (
	"log"
	"tasks/task_2/todo_api_gin/domain/usecase"

	handler "tasks/task_2/todo_api_gin/delivery"
	"tasks/task_2/todo_api_gin/domain/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	// Dependency Injection - bog'liqliklarni ulash
	todoRepo := repository.NewInMemoryTodoRepository()
	todoSvc := usecase.NewTodoUseCase(todoRepo)
	todoHandler := handler.NewTodoHandler(todoSvc)

	// Router sozlamalari
	router := gin.Default()

	// Global middleware'lar
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API versiyalash
	v1 := router.Group("/api/v1")
	todoHandler.RegisterRoutes(v1)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server ishlayapti",
		})
	})

	log.Println("Server :8080 portda ishga tushdi...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Serverni ishga tushirishda xato: %v", err)
	}
}
