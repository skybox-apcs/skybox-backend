package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/shared"
)

func FolderPermissionMiddleware(fc *controllers.FolderController, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		folderID := c.Param("folderId")
		userID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()

		// Check folder permission
		hasPermission, err := fc.CheckFolderPermission(c, folderID, userID, permission)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required permission for this folder."})
			c.Abort()
			return
		}

		c.Next()
	}
}

// FilePermissionMiddleware checks if the user has the required permission for the file's parent folder
func FilePermissionMiddleware(fc *controllers.FolderController, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("fileId")
		if fileID == "" {
			shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required.", nil)
			c.Abort()
			return
		}

		// Get the file metadata to retrieve the parent folder ID
		file, err := fc.FileService.GetFileByID(c, fileID)
		if err != nil {
			shared.RespondJson(c, http.StatusForbidden, "error", "File not found or access denied.", nil)
			c.Abort()
			return
		}

		// Check folder permission for the parent folder
		userID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
		hasPermission, err := fc.CheckFolderPermission(c, file.ParentFolderID.Hex(), userID, requiredPermission)
		if err != nil || !hasPermission {
			shared.RespondJson(c, http.StatusForbidden, "error", "Permission denied.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
