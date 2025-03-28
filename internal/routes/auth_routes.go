package routes

import (
	"skybox-backend/internal/controllers"
	"skybox-backend/internal/models"
	"skybox-backend/internal/repositories"
	"skybox-backend/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewAuthRouters sets up the routes and the corresponding handlers
func NewAuthRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create a new instance of the user repository
	ur := repositories.NewUserRepository(db, models.CollectionUsers)
	ac := &controllers.AuthController{
		AuthService: services.NewAuthService(ur),
	}

	// Create a new group for the auth routes
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", ac.RegisterHandler)
		authGroup.POST("/login", ac.LoginHandler)
	}
}
