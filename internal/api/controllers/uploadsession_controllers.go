package controllers

import (
	"fmt"
	"net/http"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type UploadSessionController struct {
	UploadSessionService services.UploadSessionService
}

func NewUploadSessionController(uss services.UploadSessionService) *UploadSessionController {
	return &UploadSessionController{
		UploadSessionService: uss,
	}
}

// GetUploadSessionHandler handles the request to get an upload session by its session ID
func (usc *UploadSessionController) GetUploadSessionHandler(c *gin.Context) {
	sessionToken := c.Param("sessionToken")
	if sessionToken == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Session Token is required")
		return
	}

	session, err := usc.UploadSessionService.GetSessionRecord(c.Request.Context(), sessionToken)
	if err != nil {
		c.Error(err)
		return
	}

	if session == nil {
		shared.ErrorJSON(c, http.StatusNotFound, "Session not found")
		return
	}

	shared.SuccessJSON(c, http.StatusOK, "Session retrieved successfully", session)
}

// AddChunkHandler handles the request to add a chunk to an upload session
func (usc *UploadSessionController) AddChunkHandler(c *gin.Context) {
	type AddChunkRequest struct {
		ChunkNumber int `json:"chunk_number" binding:"gte=0"`
	}

	var requestBody AddChunkRequest

	// Read the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		fmt.Println("Error binding JSON:", err)
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	sessionToken := c.Param("sessionToken")
	if sessionToken == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Session Token is required")
		return
	}

	err := usc.UploadSessionService.AddChunkSessionRecord(c.Request.Context(), sessionToken, requestBody.ChunkNumber)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessJSON(c, http.StatusOK, "Chunk added successfully", nil)
}
