package controllers

import (
	"net/http"
	"skybox-backend/internal/api/services"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

type SearchResult struct {
	ID     string `json:"id"`
	IsFile bool   `json:"is_file"`
	Name   string `json:"name"`
}

// SearchFilesAndFoldersHandler handles the search request
// @Summary Search files and folders
// @Description Search for files and folders by query
// @Tags Search
// @Accept json
// @Produce json
// @Param ownerId query string true "Owner ID"
// @Param query query string true "Search query"
// @Success 200 {array} SearchResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /search [get]
func (sc *SearchController) SearchFilesAndFoldersHandler(c *gin.Context) {
	ownerId := c.MustGet("x-user-id-hex").(primitive.ObjectID)
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	results, err := sc.SearchService.SearchFilesAndFolders(c.Request.Context(), ownerId, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
