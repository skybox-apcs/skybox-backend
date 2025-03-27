package routes

import (
	"skybox-backend/internal/controllers"
	"skybox-backend/internal/models"
	"skybox-backend/internal/repositories"
	"skybox-backend/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes sets up the routes and the corresponding handlers
func NewAuthRouters(db *mongo.Database, group *gin.RouterGroup) {
	ur := repositories.NewUserRepository(db, models.CollectionUsers)
	ac := &controllers.AuthController{
		AuthService: services.NewAuthService(ur),
	}

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", ac.Register)
	}
}
