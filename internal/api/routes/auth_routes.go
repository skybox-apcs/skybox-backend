package routes

import (
	"skybox-backend/configs"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewAuthRouters sets up the routes and the corresponding handlers
func NewAuthRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db)
	ac := appContainer.AuthController

	// Create a new group for the auth routes
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", ac.RegisterHandler)
		authGroup.POST("/login", ac.LoginHandler)
		authGroup.POST("/refresh", ac.RefreshHandler)
		authGroup.POST("/logout", middlewares.JwtAuthMiddleware(configs.Config.JWTSecret), ac.LogoutHandler)
	}
}
