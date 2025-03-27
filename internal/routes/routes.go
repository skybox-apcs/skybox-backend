package routes

import (
	// "skybox-backend/internal/app/controllers"

	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

// SetupRouter sets up the routes and the corresponding handlers
func SetupRouter(router *gin.Engine) *gin.Engine {
	api := router.Group("/api/v1")
	{
		api.GET("/hello", HelloWorld)
	}

	return router
}
