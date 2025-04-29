package routes

import (
	"skybox-backend/internal/blockserver/controllers"
	"skybox-backend/internal/blockserver/services"

	"github.com/gin-gonic/gin"
)

// NewUploadRouters sets up the routes and the corresponding handlers for file uploads
func NewUploadRouters(group *gin.RouterGroup) {
	// Initialize the upload service (if needed)
	uploadService := services.NewUploadService()
	// Create a new instance of the UploadController
	uploadController := controllers.NewUploadController(
		uploadService,
	)

	// Create a new group for the upload routes
	uploadGroup := group.Group("/upload")
	{
		uploadGroup.POST("/whole/:fileId", uploadController.UploadWholeFileHandler)
		uploadGroup.POST("/chunk/:fileId", uploadController.UploadAutoChunkHandler)
		uploadGroup.POST("/chunk/session/:sessionId/:chunkIndex", uploadController.UploadChunkHandler)
	}

	// Add any other upload-related routes here
}
