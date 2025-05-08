package controllers

import (
	"net/http"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type UploadSessionController struct {
	UploadSessionService *services.UploadSessionService
}

func NewUploadSessionController(uss *services.UploadSessionService) *UploadSessionController {
	return &UploadSessionController{
		UploadSessionService: uss,
	}
}

// GetUploadSessionHandler godoc
//
//	@Summary		Get an upload session by session token
//	@Description	Retrieve an upload session's metadata using its session token.
//	@Security		Bearer
//	@Tags			UploadSession
//	@Accept			json
//	@Produce		json
//	@Param			sessionToken	path		string	true	"Session Token"
//	@Success		200				{object}	models.UploadSession	"Session retrieved successfully"
//	@Failure		400				{string}	string	"Bad Request: Missing or invalid session token"
//	@Failure		404				{string}	string	"Not Found: Session not found"
//	@Failure		500				{string}	string	"Internal Server Error"
//	@Router			/api/v1/upload/{sessionToken} [get]
//
// GetUploadSessionHandler handles the request to get an upload session by its session ID
func (usc *UploadSessionController) GetUploadSessionHandler(c *gin.Context) {
	sessionToken := c.Param("sessionToken")
	if sessionToken == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Session Token is required")
		return
	}

	session, err := usc.UploadSessionService.GetSessionRecord(c, sessionToken)
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

// GetSessionRecordByFileIDHandler godoc
//
//	@Summary		Get an upload session by file ID
//	@Description	Retrieve an upload session's metadata using its associated file ID.
//	@Tags			UploadSession
//	@Accept			json
//	@Produce		json
//	@Param			fileID	path		string	true	"File ID"
//	@Success		200		{object}	models.UploadSession	"Session retrieved successfully"
//	@Failure		400		{string}	string	"Bad Request: Missing or invalid file ID"
//	@Failure		404		{string}	string	"Not Found: Session not found"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/upload/file/{fileID} [get]
//
// GetSessionRecordByFileIDHandler handles the request to get an upload session by file ID
func (usc *UploadSessionController) GetSessionRecordByFileIDHandler(c *gin.Context) {
	fileID := c.Param("fileID")
	if fileID == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "File ID is required")
		return
	}

	session, err := usc.UploadSessionService.GetSessionRecordByFileID(c, fileID)
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

// GetSessionRecordByUserIDHandler godoc
//
//	@Summary		Get all upload sessions by user ID
//	@Description	Retrieve all upload sessions associated with a user ID.
//	@Security		Bearer
//	@Tags			UploadSession
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID (me for current user)"
//	@Param			status	query		string	false	"Status of the upload session (optional)"
//	Success		200				{array}	models.UploadSession	"Sessions retrieved successfully"
//	@Failure		400				{string}	string	"Bad Request: Missing or invalid user ID"
//	@Failure		404				{string}	string	"Not Found: No sessions found for the user ID"
//	@Failure		500				{string}	string	"Internal Server Error"
//	@Router			/upload/user/{userID} [get]
//
// GetSessionRecordByUserIDHandler handles the request to get all upload sessions by user ID
func (usc *UploadSessionController) GetSessionRecordByUserIDHandler(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" || userID == "me" {
		// Default to the user ID from the JWT token if not provided in the URL
		userID = c.GetString("x-user-id")

		if userID == "" {
			shared.ErrorJSON(c, http.StatusBadRequest, "User ID is required")
			return
		}
	}

	statusQuery := c.Query("status")
	if statusQuery != "" {
		if statusQuery != "pending" && statusQuery != "completed" {
			shared.ErrorJSON(c, http.StatusBadRequest, "Invalid status query parameter")
			return
		}
	}

	sessions, err := usc.UploadSessionService.GetSessionRecordByUserID(c, userID)
	if err != nil {
		c.Error(err)
		return
	}

	if sessions == nil || len(*sessions) == 0 {
		shared.ErrorJSON(c, http.StatusNotFound, "No sessions found for the user ID")
		return
	}

	// Filter sessions based on status if provided
	if statusQuery != "" {
		var filteredSessions []models.UploadSession
		for _, session := range *sessions {
			if session.Status == statusQuery {
				filteredSessions = append(filteredSessions, session)
			}
		}
		sessions = &filteredSessions
	}

	shared.SuccessJSON(c, http.StatusOK, "Sessions retrieved successfully", sessions)
}

// AddChunkHandler godoc
//
//	@Summary		Add a chunk to an upload session
//	@Description	Add a chunk to an existing upload session using its session token.
//	@Tags			UploadSession
//	@Accept			json
//	@Produce		json
//	@Param			sessionToken	path		string	true	"Session Token"
//	@Param			body			body		models.AddChunkRequest	true	"Chunk data"
//	@Success		200				{string}	string	"Chunk added successfully"
//	@Failure		400				{string}	string	"Bad Request: Invalid request body or session token"
//	@Failure		500				{string}	string	"Internal Server Error"
//	@Router			/upload/{sessionToken} [put]
//
// AddChunkHandler handles the request to add a chunk to an upload session
func (usc *UploadSessionController) AddChunkHandler(c *gin.Context) {
	var requestBody models.AddChunkRequest

	// Read the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	sessionToken := c.Param("sessionToken")
	if sessionToken == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Session Token is required")
		return
	}

	err := usc.UploadSessionService.AddChunkSessionRecord(c, sessionToken, requestBody.ChunkNumber, requestBody.ChunkSize, requestBody.ChunkHash)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessJSON(c, http.StatusOK, "Chunk added successfully", nil)
}

// AddChunkViaFileIDHandler godoc
//
//	@Summary		Add a chunk to an upload session using file ID
//	@Description	Add a chunk to an existing upload session using its associated file ID.
//	@Tags			UploadSession
//	@Accept			json
//	@Produce		json
//	@Param			fileID	path		string	true	"File ID"
//	@Param			body	body		models.AddChunkViaFileIDRequest	true	"Chunk data"
//	@Success		200		{string}	string	"Chunk added successfully"
//	@Failure		400		{string}	string	"Bad Request: Invalid request body or file ID"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/upload/file/{fileID} [put]
//
// AddChunkViaFileIDHandler handles the request to add a chunk to an upload session using file ID
func (usc *UploadSessionController) AddChunkViaFileIDHandler(c *gin.Context) {
	var requestBody models.AddChunkViaFileIDRequest

	// Read the request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	fileID := c.Param("fileID")
	if fileID == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "File ID is required")
		return
	}

	err := usc.UploadSessionService.AddChunkSessionRecordByFileID(c, fileID, requestBody.ChunkNumber, requestBody.ChunkSize, requestBody.ChunkHash)
	if err != nil {
		c.Error(err)
		return
	}

	shared.SuccessJSON(c, http.StatusOK, "Chunk added successfully", nil)
}
