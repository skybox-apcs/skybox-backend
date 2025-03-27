package routes

import (
	"github.com/gin-gonic/gin"
)

func HelloWorldFunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

// SetupRouter sets up the routes and the corresponding handlers
func SetupRouter(router *gin.Engine) *gin.Engine {
	api := router.Group("/api/v1")
	{
		api.GET("/hello", HelloWorldFunc)
	}

	// Add the auth routes
	RegisterLoginRoutes(api)

	return router
}
