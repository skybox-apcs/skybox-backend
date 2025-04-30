package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUploadRouters sets up the routes and the corresponding handlers
func NewUploadRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create new instance of the repositories
	uploadSessionRepo := repositories.NewUploadSessionRepository(db, models.CollectionUploadSessions)
	usc := controllers.NewUploadSessionController(
		*services.NewUploadSessionService(uploadSessionRepo),
	)

	// Create a new group for the folder routes
	folderGroup := group.Group("/upload")
	{
		folderGroup.GET("/:sessionToken", usc.GetUploadSessionHandler)
		folderGroup.PUT("/:sessionToken", usc.AddChunkHandler)
		folderGroup.GET("/file/:fileID", usc.GetSessionRecordByFileIDHandler)
		folderGroup.PUT("/file/:fileID", usc.AddChunkViaFileIDHandler)
	}
}
