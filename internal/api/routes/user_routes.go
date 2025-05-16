package routes

import (
	"skybox-backend/configs"
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserRouters sets up the routes and the corresponding handlers
func NewUserRouters(db *mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, models.CollectionUsers)
	uc := &controllers.UserController{
		UserService: services.NewUserService(ur),
	}

	userGroup := group.Group("/user")
	publicGroup := userGroup.Group("")
	// Public Routes
	{
		publicGroup.GET("/:userId", uc.GetUserByIDHandler)
		publicGroup.GET("/email/:email", uc.GetUserByEmailHandler)
		publicGroup.POST("/ids", uc.GetUsersByIDsHandler)
		publicGroup.GET("/emails", uc.GetUsersByEmailsHandler)
	}

	privateGroup := userGroup.Group("")
	privateGroup.Use(middlewares.JwtAuthMiddleware(configs.Config.JWTSecret))
	// Private Routes
	{
		userGroup.GET("/info", uc.GetUserInformationHandler)
	}
}
