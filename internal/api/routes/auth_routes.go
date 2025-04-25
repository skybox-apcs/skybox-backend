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

// NewAuthRouters sets up the routes and the corresponding handlers
func NewAuthRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create a new instance of the repositories
	ur := repositories.NewUserRepository(db, models.CollectionUsers)
	utr := repositories.NewUserTokenRepository(db, models.CollectionUserTokens)
	ac := &controllers.AuthController{
		AuthService:      services.NewAuthService(ur),
		UserTokenService: services.NewUserTokenService(utr),
	}

	// Create a new group for the auth routes
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", ac.RegisterHandler)
		authGroup.POST("/login", ac.LoginHandler)
		authGroup.POST("/refresh", ac.RefreshHandler)
		authGroup.POST("/logout", middlewares.JwtAuthMiddleware(configs.Config.JWTSecret), ac.LogoutHandler)
	}
}
