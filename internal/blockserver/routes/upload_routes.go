package routes

import (
	"skybox-backend/internal/blockserver/controllers"

	"github.com/gin-gonic/gin"
)

// NewUploadRouters sets up the routes and the corresponding handlers for file uploads
func NewUploadRouters(group *gin.RouterGroup) {
	// Create a new instance of the UploadController
	uploadController := controllers.NewUploadController()

	// Create a new group for the upload routes
	uploadGroup := group.Group("/upload")
	{
		uploadGroup.POST("/whole", uploadController.UploadWholeFileHandler)
		uploadGroup.POST("/chunk/:sessionId/:chunkIndex", uploadController.UploadChunkHandler)
	}

	// Add any other upload-related routes here
}
