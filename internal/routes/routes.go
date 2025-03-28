package routes

import (
	"skybox-backend/configs"
	"skybox-backend/internal/controllers"
	"skybox-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes sets up the routes and the corresponding handlers
func SetupRouter(db *mongo.Database, gin *gin.Engine) *gin.Engine {
	// Swagger routes
	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRouter := gin.Group("")

	// Setup the v1 routes
	v1 := publicRouter.Group("/api/v1")

	// Public routes
	{
		// Setup the auth routes
		NewAuthRouters(db, v1)

		// Hello World routes
		v1.GET("/hello", controllers.HelloWorldHandler)
	}

	// Private routes
	protectedRouter := gin.Group("")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(configs.Config.JWTSecret))

	v1 = protectedRouter.Group("/api/v1")

	{
		// Setup the user routes
		NewUserRouters(db, v1)
	}

	return gin
}
