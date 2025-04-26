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
	fr := repositories.NewFolderRepository(db, models.CollectionFolders)
	fc := &controllers.FolderController{
		FolderService: services.NewFolderService(fr),
	}

	// Create a new group for the folder routes
	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId", fc.GetFolderHandler)
		folderGroup.GET("/:folderId/contents", fc.GetContentsHandler)
		folderGroup.POST("/:folderId/create", fc.CreateFolderHandler)
		folderGroup.DELETE("/:folderId", fc.DeleteFolderHandler)
	}
}
