package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewFolderRouters sets up the routes and the corresponding handlers
func NewFolderRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create new instance of the repositories
	folderRepo := repositories.NewFolderRepository(db, models.CollectionFolders)
	fileRepo := repositories.NewFileRepository(db, models.CollectionFiles)
	fc := &controllers.FolderController{
		FolderService: services.NewFolderService(folderRepo),
		FileService:   services.NewFileService(fileRepo),
	}

	// Create a new group for the folder routes
	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId", fc.GetFolderHandler)
		folderGroup.DELETE("/:folderId", fc.DeleteFolderHandler)
		folderGroup.GET("/:folderId/contents", fc.GetContentsHandler)
		folderGroup.POST("/:folderId/create", fc.CreateFolderHandler)
		folderGroup.PUT("/:folderId/rename", fc.RenameFolderHandler)
		folderGroup.PATCH("/:folderId/rename", fc.RenameFolderHandler)
		folderGroup.PUT("/:folderId/move", fc.MoveFolderHandler)

		folderGroup.POST("/:folderId/upload", fc.UploadFileMetadataHandler) // TODO: Implement upload file metadata handler
	}
}
