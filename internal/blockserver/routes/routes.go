package routes

import (
	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/controllers"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes and the corresponding handlers
func SetupRouter(gin *gin.Engine) *gin.Engine {
	// Swagger routes
	// gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRouter := gin.Group("")

	// Setup the v1 routes
	v1 := publicRouter.Group("")

	// Public routes
	{
		// Hello World routes
		v1.GET("/hello", controllers.HelloWorldHandler)

		// Download routes
		NewDownloadRouters(v1)
	}

	// Private routes
	protectedRouter := gin.Group("")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(configs.Config.JWTSecret))

	v1 = protectedRouter.Group("")

	{
		v1.GET("/protected/hello", controllers.HelloWorldHandler)

		// Upload routes
		NewUploadRouters(v1)
	}

	return gin
}
