package routes

import (
	"skybox-backend/internal/blockserver/controllers"
	"skybox-backend/internal/blockserver/services"

	"github.com/gin-gonic/gin"
)

// NewDownloadRouters sets up the routes and the corresponding handlers for file downloads
func NewDownloadRouters(group *gin.RouterGroup) {
	// Initialize the download service (if needed)
	downloadService := services.NewDownloadService()
	// Create a new instance of the DownloadController
	downloadController := controllers.NewDownloadController(
		downloadService,
	)

	// Create a new group for the download routes
	downloadGroup := group.Group("/download")
	{
		// Download a file by its ID
		downloadGroup.GET("/:fileId", downloadController.DownloadFileHandler)
	}
}
