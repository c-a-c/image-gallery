package router

import (
	"backend/controller"
	"backend/infrastructure/middleware"
	"backend/repository"
	"backend/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupRoutes configures all routes
func SetupRoutes(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.CORSMiddleware())
	e.Use(middleware.LoggerMiddleware())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	imageRepo := repository.NewImageRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	imageUseCase := usecase.NewImageUseCase(imageRepo)
	postUseCase := usecase.NewPostUseCase(postRepo, imageRepo)

	// Initialize controllers
	authController := controller.NewAuthController(userUseCase)
	imageController := controller.NewImageController(imageUseCase)
	postController := controller.NewPostController(postUseCase)

	// Public routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Image Gallery API",
			"version": "1.0.0",
		})
	})

	// Auth routes
	auth := e.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// User routes
	api.GET("/profile", authController.GetProfile)
	api.PUT("/profile", authController.UpdateProfile)
	api.PUT("/password", authController.ChangePassword)

	// Image routes
	api.POST("/images", imageController.UploadImage)
	api.GET("/images/my", imageController.GetUserImages)
	api.GET("/images/:id", imageController.GetImage)
	api.PUT("/images/:id", imageController.UpdateImage)
	api.DELETE("/images/:id", imageController.DeleteImage)

	// Post routes
	api.POST("/posts", postController.CreatePost)
	api.GET("/posts/my", postController.GetUserPosts)
	api.GET("/posts/:id", postController.GetPost)
	api.PUT("/posts/:id", postController.UpdatePost)
	api.DELETE("/posts/:id", postController.DeletePost)

	// Public routes (no auth required)
	public := e.Group("/public")
	public.Use(middleware.OptionalAuthMiddleware())

	// Public image routes
	public.GET("/images", imageController.GetPublicImages)
	public.GET("/images/search", imageController.SearchImages)

	// Public post routes
	public.GET("/posts", postController.GetPublicPosts)
	public.GET("/posts/search", postController.SearchPosts)
	public.GET("/posts/tags", postController.GetPostsByTags)

	return e
}
