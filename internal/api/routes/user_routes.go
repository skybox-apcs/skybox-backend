package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"

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
	{
		userGroup.GET("/info", uc.GetUserInformationHandler)
	}
}
