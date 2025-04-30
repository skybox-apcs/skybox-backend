package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewFileRouters sets up the routes and the corresponding handlers
func NewFileRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db)
	fc := appContainer.FileController

	folderRepo := repositories.NewFolderRepository(db, models.CollectionFolders)
	folderController := &controllers.FolderController{
		FolderService: services.NewFolderService(folderRepo),
		FileService:   services.NewFileService(fr),
	}

	// Create a new group for the file routes
	fileGroup := group.Group("/files")
	{
		fileGroup.GET("/:fileId", fc.GetFileMetadataHandler)
		fileGroup.DELETE("/:fileId", fc.DeleteFileHandler)
		fileGroup.PUT("/:fileId/rename", fc.RenameFileHandler)
		fileGroup.PATCH("/:fileId/rename", fc.RenameFileHandler)
		fileGroup.PUT("/:fileId/move", fc.MoveFileHandler)

		fileGroup.GET("/:fileId/download", fc.FullDownloadFileHandler)
		fileGroup.GET("/:fileId/partial_download", fc.PartialDownloadFileHandler)
	}
}
