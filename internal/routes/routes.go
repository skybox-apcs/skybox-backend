package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

// SetupRoutes sets up the routes and the corresponding handlers
func SetupRouter(db *mongo.Database, gin *gin.Engine) *gin.Engine {
	publicRouter := gin.Group("")

	// Setup the v1 routes
	v1 := publicRouter.Group("/api/v1")

	// Public routes
	{
		// Setup the auth routes
		NewAuthRouters(db, v1)

		v1.GET("/", HelloWorld)
	}
	// Private routes
	{
		// ...
	}

	return gin
}
