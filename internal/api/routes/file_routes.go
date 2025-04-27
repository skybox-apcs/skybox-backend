package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewFileRouters sets up the routes and the corresponding handlers
func NewFileRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create new instance of the repositories
	fr := repositories.NewFileRepository(db, models.CollectionFiles)
	fc := &controllers.FileController{
		FileService: services.NewFileService(fr),
	}

	// Create a new group for the file routes
	fileGroup := group.Group("/files")
	{
		fileGroup.GET("/:fileId", fc.GetFileMetadataHandler)
		fileGroup.DELETE("/:fileId", fc.DeleteFileHandler)
		fileGroup.PUT("/:fileId/rename", fc.RenameFileHandler)
		fileGroup.PATCH("/:fileId/rename", fc.RenameFileHandler)
		fileGroup.PUT("/:fileId/move", fc.MoveFileHandler)
	}
}
