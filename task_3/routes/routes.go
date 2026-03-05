package routes

import (
	"net/http"
	"tasks/task_3/config"
	"tasks/task_3/handlers"
	"tasks/task_3/middleware"
	"tasks/task_3/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	authRepo := repository.NewAuthRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	r.Use(middleware.CORSMiddleware())

	authHandler := handlers.NewAuthHandler(cfg.JWTSecret, authRepo)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(commentRepo)
	userHandler := handlers.NewUserHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	tagHandler := handlers.NewTagHandler(db)
	searchHandler := handlers.NewSearchHandler(db)
	statsHandler := handlers.NewStatsHandler(db)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Blog API ishlayapti! 🚀"})
	})

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	authRequired := api.Group("/")
	authRequired.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authRequired.GET("/auth/me", authHandler.GetMe)
		authRequired.PUT("/auth/profile", authHandler.UpdateProfile)
		authRequired.PUT("/auth/change-password", authHandler.ChangePassword)
	}

	posts := api.Group("/posts")
	{
		// Ochiq route'lar
		posts.GET("", postHandler.GetPosts)                 // GET /api/v1/posts
		posts.GET("/:id", postHandler.GetPost)              // GET /api/v1/posts/1
		posts.GET("/slug/:slug", postHandler.GetPostBySlug) // GET /api/v1/posts/slug/my-post

		// Kommentariylar (ochiq - ko'rish)
		posts.GET("/:id/comments", commentHandler.GetComments) // GET /api/v1/posts/1/comments

		// Auth kerak
		postsAuth := posts.Group("/")
		postsAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			postsAuth.POST("", postHandler.CreatePost)                    // POST /api/v1/posts
			postsAuth.PUT("/:id", postHandler.UpdatePost)                 // PUT /api/v1/posts/1
			postsAuth.DELETE("/:id", postHandler.DeletePost)              // DELETE /api/v1/posts/1
			postsAuth.POST("/:id/like", postHandler.LikePost)             // POST /api/v1/posts/1/like
			postsAuth.POST("/:id/comments", commentHandler.CreateComment) // POST /api/v1/posts/1/comments
		}
	}

	comments := api.Group("/comments")
	comments.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		comments.PUT("/:id", commentHandler.UpdateComment)     // PUT /api/v1/comments/1
		comments.DELETE("/:id", commentHandler.DeleteComment)  // DELETE /api/v1/comments/1
		comments.POST("/:id/like", commentHandler.LikeComment) // POST /api/v1/comments/1/like
	}

	users := api.Group("/users")
	{
		// Ochiq
		users.GET("/:id", userHandler.GetUser)                // GET /api/v1/users/1
		users.GET("/:id/posts", postHandler.GetUserPosts)     // GET /api/v1/users/1/posts
		users.GET("/:id/followers", userHandler.GetFollowers) // GET /api/v1/users/1/followers
		users.GET("/:id/following", userHandler.GetFollowing) // GET /api/v1/users/1/following

		// Auth kerak
		usersAuth := users.Group("/")
		usersAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			usersAuth.POST("/:id/follow", userHandler.FollowUser) // POST /api/v1/users/1/follow
		}
	}

	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.GetCategories)   // GET /api/v1/categories
		categories.GET("/:id", categoryHandler.GetCategory) // GET /api/v1/categories/1
	}

	tags := api.Group("/tags")
	{
		tags.GET("", tagHandler.GetTags)               // GET /api/v1/tags
		tags.GET("/:id/posts", tagHandler.GetTagPosts) // GET /api/v1/tags/1/posts

		tagsAuth := tags.Group("/")
		tagsAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			tagsAuth.POST("", tagHandler.CreateTag) // POST /api/v1/tags
		}
	}

	api.GET("/search", searchHandler.Search) // GET /api/v1/search?q=golang

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/stats", statsHandler.GetStats)                      // GET /api/v1/admin/stats
		admin.GET("/users", userHandler.GetUsers)                       // GET /api/v1/admin/users
		admin.POST("/categories", categoryHandler.CreateCategory)       // POST /api/v1/admin/categories
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory) // DELETE /api/v1/admin/categories/1
	}

	return r
}
