package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserRouters sets up the routes and the corresponding handlers
func NewUserRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db)
	uc := appContainer.UserController

	// Create a new group for the user routes
	userGroup := group.Group("/user")
	publicGroup := userGroup.Group("")
	// Public Routes
	{
		publicGroup.GET("/:userId", uc.GetUserByIDHandler)
		publicGroup.GET("/email/:email", uc.GetUserByEmailHandler)
		publicGroup.POST("/ids", uc.GetUsersByIDsHandler)
		publicGroup.POST("/emails", uc.GetUsersByEmailsHandler)
	}

	privateGroup := userGroup.Group("")
	privateGroup.Use(middlewares.JwtAuthMiddleware(configs.Config.JWTSecret))
	// Private Routes
	{
		userGroup.GET("/info", uc.GetUserInformationHandler)
	}
}
