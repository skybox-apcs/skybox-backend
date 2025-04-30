package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewFolderRouters sets up the routes and the corresponding handlers
func NewFolderRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db)
	fc := appContainer.FolderController

	// Create a new group for the folder routes
	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId", middlewares.FolderPermissionMiddleware(fc, "view"), fc.GetFolderHandler)
		folderGroup.DELETE("/:folderId", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.DeleteFolderHandler)
		folderGroup.GET("/:folderId/contents", middlewares.FolderPermissionMiddleware(fc, "view"), fc.GetContentsHandler)
		folderGroup.POST("/:folderId/create", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.CreateFolderHandler)
		folderGroup.PUT("/:folderId/rename", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.RenameFolderHandler)
		folderGroup.PATCH("/:folderId/rename", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.RenameFolderHandler)
		folderGroup.PUT("/:folderId/move", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.MoveFolderHandler)

		folderGroup.POST("/:folderId/upload", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.UploadFileMetadataHandler) // TODO: Implement upload file metadata handler

		// Share
		folderGroup.PUT("/:folderId/public-status", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.UpdateFolderPublicStatusHandler)
		folderGroup.PUT("/:folderId/public-status/all", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.UpdateFolderAndSubfoldersPublicStatusHandler)
		folderGroup.GET("/:folderId/public-status", middlewares.FolderPermissionMiddleware(fc, "view"), fc.GetFolderPublicStatusHandler)
		folderGroup.POST("/:folderId/share", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.ShareFolderHandler)
		folderGroup.DELETE("/:folderId/share", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.RemoveFolderShareHandler)
		folderGroup.GET("/:folderId/shared-users", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.GetFolderSharedUsersHandler)
		folderGroup.POST("/:folderId/share/all", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.ShareFolderAndSubfoldersHandler)
		folderGroup.DELETE("/:folderId/share/all", middlewares.FolderPermissionMiddleware(fc, "edit"), fc.RevokeFolderAndSubfoldersShareHandler)
	}
}
