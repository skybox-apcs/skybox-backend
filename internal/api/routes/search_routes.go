package routes

import (
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewSearchRouters sets up the search routes
func NewSearchRouters(db *mongo.Database, group *gin.RouterGroup) {
	// Create repositories
	fileRepo := repositories.NewFileRepository(db, models.CollectionFiles)
	folderRepo := repositories.NewFolderRepository(db, models.CollectionFolders)
	userRepo := repositories.NewUserRepository(db, models.CollectionUsers)

	// Create services
	searchService := services.NewSearchService(fileRepo, folderRepo, userRepo)

	// Create controller
	searchController := controllers.NewSearchController(searchService)

	// Add search route
	group.GET("/search", searchController.SearchFilesAndFoldersHandler)
}
