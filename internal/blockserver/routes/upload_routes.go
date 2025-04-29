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
		// Non-resumable upload
		uploadGroup.POST("/whole/:fileId", uploadController.UploadWholeFileHandler)
		uploadGroup.POST("/chunked/:fileId", uploadController.UploadAutoChunkHandler)

		// Resumable upload
		// 1. Start a new session
		// uploadGroup.POST("/session/start", nil)
		// 2. Upload a chunk to a resumable session
		uploadGroup.POST("/session/:sessionToken/chunk", uploadController.UploadChunkHandler)
		// 3. Get the status of a resumable session
		uploadGroup.GET("/session/:sessionToken/status", nil)

		// (Optional) Cancel session, merge chunks, etc.
		uploadGroup.POST("/session/:sessionToken/complete", nil)
		uploadGroup.DELETE("/session/:sessionToken", nil)

		uploadGroup.POST("/session/:sessionToken/:chunkIndex", uploadController.UploadChunkHandler)
	}

	// Add any other upload-related routes here
}
