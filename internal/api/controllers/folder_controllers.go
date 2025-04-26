package controllers

import (
	"fmt"
	"net/http"
	"time"

	// "skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderController struct {
	FolderService *services.FolderService
}

func NewFolderController(folderService *services.FolderService) *FolderController {
	return &FolderController{
		FolderService: folderService,
	}
}

func (fc *FolderController) CreateFolder(c *gin.Context) {
	// Defune the request body structure
	var request models.CreateFolderRequest

	// Bind the request body to the structure
	err := c.ShouldBindJSON(&request)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	// Get the folder parent folder ID from the request
	var parentFolderIdHex primitive.ObjectID = primitive.NilObjectID
	folderId := c.Param("folderId")
	if folderId != "" {
		var parentFolderId string = ""
		parentFolderId, err = fc.FolderService.GetFolderParentIDByFolderID(c, folderId)
		if err != nil {
			shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to get parent folder ID.: Error: "+err.Error(), nil)
			return
		}

		// Cast the parent folder ID to ObjectID
		parentFolderIdHex, err = primitive.ObjectIDFromHex(parentFolderId)
		if err != nil {
			shared.RespondJson(c, http.StatusInternalServerError, "error", "Invalid parent folder ID. Error: "+err.Error(), nil)
			return
		}
	}

	// Cast the owner ID to ObjectID
	// c.Set("x-user-id", userId)
	ownerIdHex, err := primitive.ObjectIDFromHex(c.GetString("x-user-id"))
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid owner ID.", nil)
		return
	}

	// Create the folder object
	folder := &models.Folder{
		Name:           request.Name,
		OwnerID:        ownerIdHex,
		ParentFolderID: parentFolderIdHex,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      nil,
	}

	// Create the folder in the database
	folderResult, err := fc.FolderService.CreateFolder(c, folder)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to create folder.", nil)
		return
	}

	// Create the response object
	response := models.CreateFolderResponse{
		ID:             folderResult.ID.Hex(),
		OwnerID:        folder.OwnerID.Hex(),
		ParentFolderID: folder.ParentFolderID.Hex(),
		Name:           folder.Name,
	}

	// Send the response
	shared.RespondJson(c, http.StatusCreated, "success", "Folder created successfully.", response)
}

// GetContents retrieves the contents of a folder
func (fc *FolderController) GetContents(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")

	// Get the folder contents from the service
	contents, err := fc.FolderService.GetFolderContents(c, folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to get folder contents. Error: "+err.Error(), nil)
		return
	}

	fmt.Println("Folder contents retrieved successfully:", contents)

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "Folder contents retrieved successfully.", contents)
}
