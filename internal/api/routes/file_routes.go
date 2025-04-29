package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewFileRouters sets up the routes and the corresponding handlers
func NewFileRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create new instance of the repositories
	fr := repositories.NewFileRepository(db, models.CollectionFiles)
	usr := repositories.NewUploadSessionRepository(db, models.CollectionUploadSessions)
	fc := &controllers.FileController{
		FileService: services.NewFileService(fr, usr),
	}

	folderRepo := repositories.NewFolderRepository(db, models.CollectionFolders)
	folderController := &controllers.FolderController{
		FolderService: services.NewFolderService(folderRepo),
		FileService:   services.NewFileService(fr),
	}

	// Create a new group for the file routes
	fileGroup := group.Group("/files")
	{
		fileGroup.GET("/:fileId", middlewares.FilePermissionMiddleware(folderController, "view"), fc.GetFileMetadataHandler)
		fileGroup.DELETE("/:fileId", middlewares.FilePermissionMiddleware(folderController, "edit"), fc.DeleteFileHandler)
		fileGroup.PUT("/:fileId/rename", middlewares.FilePermissionMiddleware(folderController, "edit"), fc.RenameFileHandler)
		fileGroup.PATCH("/:fileId/rename", middlewares.FilePermissionMiddleware(folderController, "edit"), fc.RenameFileHandler)
		fileGroup.PUT("/:fileId/move", middlewares.FilePermissionMiddleware(folderController, "edit"), fc.MoveFileHandler)
	}
}
