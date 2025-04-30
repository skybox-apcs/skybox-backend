package controllers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"skybox-backend/internal/api/services"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	SearchService *services.SearchService
}

// NewSearchController creates a new instance of SearchController
func NewSearchController(searchService *services.SearchService) *SearchController {
	return &SearchController{
		SearchService: searchService,
	}
}

// SearchFilesAndFoldersHandler handles the search request
// @Summary Search files and folders
// @Description Search for files and folders by query
// @Tags Search
// @Accept json
// @Produce json
// @Param ownerId query string true "Owner ID"
// @Param query query string true "Search query"
// @Success 200 {array} struct{ID string `json:"id"`; IsFile bool `json:"is_file"`; Name string `json:"name"`}
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /search [get]
func (sc *SearchController) SearchFilesAndFoldersHandler(c *gin.Context) {
	ownerId := c.MustGet("x-user-id-hex").(primitive.ObjectID)
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query are required"})
		return
	}

	results, err := sc.SearchService.SearchFilesAndFolders(c.Request.Context(), ownerId, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
