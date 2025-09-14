package routes

import (
	"compass-backend/config"
	"compass-backend/internal/controllers"
	"compass-backend/internal/middleware"
	"compass-backend/internal/repositories"
	"compass-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	specRepo := repositories.NewSpecificationRepository(db)
	rfiRepo := repositories.NewRFIRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo, specRepo, rfiRepo)
	specService := services.NewSpecificationService(specRepo)
	rfiService := services.NewRFIService(rfiRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	projectController := controllers.NewProjectController(projectService)
	specController := controllers.NewSpecificationController(specService)
	rfiController := controllers.NewRFIController(rfiService)

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/signin", authController.SignIn)
		auth.POST("/refresh", authController.RefreshToken)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		// Logout
		api.POST("/auth/logout", authController.Logout)

		// User management (Admin only)
		users := api.Group("/users")
		users.Use(middleware.AdminOnly())
		{
			users.POST("", userController.CreateUser)
			users.PATCH("/:id/status", userController.UpdateUserStatus)
			users.GET("", userController.ListUsers)
		}

		// Projects
		projects := api.Group("/projects")
		{
			projects.POST("", projectController.CreateProject)
			projects.GET("", projectController.ListProjects)
			projects.GET("/:id", projectController.GetProject)
			projects.PATCH("/:id/status", projectController.UpdateProjectStatus)
			projects.DELETE("/:id", middleware.AdminOnly(), projectController.DeleteProject)

			// Project specifications
			projects.POST("/:id/specifications", specController.CreateSpecification)
			projects.GET("/:id/specifications", specController.GetProjectSpecifications)

			// Project RFIs
			projects.POST("/:id/rfis", rfiController.CreateRFI)
			projects.GET("/:id/rfis", rfiController.GetProjectRFIs)
		}

		// RFIs
		rfis := api.Group("/rfis")
		{
			rfis.PATCH("/:id/answer", rfiController.AnswerRFI)
		}
	}
}