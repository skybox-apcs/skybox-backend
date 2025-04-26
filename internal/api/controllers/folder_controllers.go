package controllers

import (
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

// CreateFolderHandler godoc
//
// @Summary Create a new folder in prompted folder id
// @Description Create a new folder in the specified parent folder
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.CreateFolderRequest true "Create Folder Request"
// @Success 201 {object} models.CreateFolderResponse
// @Failure 400 {string} string "Invalid request.
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId}/create [post]
func (fc *FolderController) CreateFolderHandler(c *gin.Context) {
	// Defune the request body structure
	var request models.CreateFolderRequest

	// Bind the request body to the structure
	err := c.ShouldBindJSON(&request)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Cast the parent folder ID to ObjectID
	parentFolderIdHex, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Invalid parent folder ID. Error: "+err.Error(), nil)
		return
	}

	// Cast the owner ID to ObjectID
	// c.Set("x-user-id", userId)
	ownerIdHex, err := primitive.ObjectIDFromHex(c.GetString("x-user-id"))
	if err != nil {
		shared.RespondJson(c, http.StatusNotFound, "error", "Invalid owner ID.", nil)
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

// GetFolderContentsHandler godoc
//
// @Summary Get contents of a folder (folders and files)
// @Description Retrieve all folders and files inside the specified parent folder
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {object} models.GetFolderContentsResponse
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId}/contents [get]
func (fc *FolderController) GetContentsHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Get the folder list from the service
	folderList, err := fc.FolderService.GetFolderListInFolder(c, folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to get folder contents. Error: "+err.Error(), nil)
		return
	}

	// Get the file list from the service
	fileList, err := fc.FolderService.GetFileListInFolder(c, folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to get file list. Error: "+err.Error(), nil)
		return
	}

	// Check if folderList and fileList are nil and initialize them to empty slices
	if folderList == nil {
		folderList = []*models.Folder{}
	}
	if fileList == nil {
		fileList = []*models.File{}
	}

	// Create the response object
	contents := models.GetFolderContentsResponse{
		FolderList: folderList,
		FileList:   fileList,
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "Folder contents retrieved successfully.", contents)
}

// DeleteFolderHandler godoc
func (fc *FolderController) DeleteFolderHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Delete the folder using the service
	err := fc.FolderService.DeleteFolder(c, folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to delete folder. Error: "+err.Error(), nil)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder deleted successfully.", nil)
}
