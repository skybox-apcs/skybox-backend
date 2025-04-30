package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUploadRouters sets up the routes and the corresponding handlers
func NewUploadRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db)
	usc := appContainer.UploadSessionController

	// Create a new group for the upload routes
	folderGroup := group.Group("/upload")
	{
		folderGroup.GET("/:sessionToken", usc.GetUploadSessionHandler)
		folderGroup.PUT("/:sessionToken", usc.AddChunkHandler)
		folderGroup.GET("/file/:fileID", usc.GetSessionRecordByFileIDHandler)
		folderGroup.PUT("/file/:fileID", usc.AddChunkViaFileIDHandler)
		folderGroup.GET("/user/:userID", usc.GetSessionRecordByUserIDHandler)
	}
}
