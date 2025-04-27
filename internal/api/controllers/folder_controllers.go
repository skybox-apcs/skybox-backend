package controllers

import (
	"net/http"
	"path/filepath"
	"sync"
	"time"

	// "skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FolderController struct {
	FileService   *services.FileService
	FolderService *services.FolderService
}

func NewFolderController(folderService *services.FolderService, fileService *services.FileService) *FolderController {
	return &FolderController{
		FolderService: folderService,
		FileService:   fileService,
	}
}

// GetFolderHandler godoc
//
// @Summary Get metadata of a specific folder
// @Description Retrieve detailed metadata for a folder by its ID. Includes information such as folder ID, owner ID, name, parent folder ID, and creation timestamps.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {object} models.Folder
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId} [get]
func (fc *FolderController) GetFolderHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	folder, err := fc.FolderService.GetFolderByID(c, folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to get folder. Error: "+err.Error(), nil)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "Folder retrieved successfully.", folder)
}

// CreateFolderHandler godoc
//
// @Summary Create a new folder inside a specified parent folder
// @Description Create a new folder inside a given parent folder. If no parent folder ID is provided, the folder will be created at the root level.
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
	ownerIdHex := c.MustGet("x-user-id-hex").(primitive.ObjectID)

	// Create the folder object
	folder := &models.Folder{
		Name:           request.Name,
		OwnerID:        ownerIdHex,
		ParentFolderID: parentFolderIdHex,
		IsDeleted:      false,
		IsRoot:         false,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
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
// @Summary List all files and folders inside a folder
// @Description Retrieve a list of all files and subfolders contained within a specified folder. Useful for browsing the contents of a directory.
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
	var folderList []*models.FolderResponse
	var fileList []*models.FileResponse
	var folderErr, fileErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		folderList, folderErr = fc.FolderService.GetFolderResponseListInFolder(c, folderId)
	}()

	go func() {
		defer wg.Done()
		fileList, fileErr = fc.FolderService.GetFileResponseListInFolder(c, folderId)
	}()

	wg.Wait()

	if folderErr != nil || fileErr != nil {
		shared.RespondJson(c, http.StatusNotFound, "error", "Failed to get folder contents. Error: "+folderErr.Error(), nil)
		return
	}

	// Check if folderList and fileList are nil and initialize them to empty slices
	if folderList == nil {
		folderList = []*models.FolderResponse{}
	}
	if fileList == nil {
		fileList = []*models.FileResponse{}
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
//
// @Summary Soft-delete a folder
// @Description Mark a folder as deleted (move it to the recycle bin). Subfolders and files are not immediately deleted but will also be marked for deletion. After a retention period, they may be permanently deleted.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {string} string "Folder deleted successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId} [delete]
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
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to delete folder.", nil)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder deleted successfully.", nil)
}

// RenameFolderHandler godoc
//
// @Summary Rename a specific folder
// @Description Rename an existing folder by providing its ID and the new name. Only the folder name will be updated; the folder's contents are unaffected.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.RenameFolderRequest true "Rename Folder Request"
// @Success 200 {string} string "Folder renamed successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId}/rename [put]
// @Router /folders/{folderId}/rename [patch]
func (fc *FolderController) RenameFolderHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Define the request body structure
	var request models.RenameFolderRequest

	// Bind the request body to the structure
	err := c.ShouldBindJSON(&request)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	// Rename the folder using the service
	err = fc.FolderService.RenameFolder(c, folderId, request.NewName)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to rename folder.", nil)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder renamed successfully.", nil)
}

// MoveFolderHandler godoc
//
// @Summary Move a folder to a new parent folder
// @Description Move a folder from its current location to another destination folder. The folder's contents will be moved along with it.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.MoveFolderRequest true "Move Folder Request"
// @Success 200 {string} string "Folder moved successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId}/move [put]
func (fc *FolderController) MoveFolderHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Define the request body structure
	var request models.MoveFolderRequest

	// Bind the request body to the structure
	err := c.ShouldBindJSON(&request)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	// Move the folder using the service
	err = fc.FolderService.MoveFolder(c, folderId, request.NewParentID)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to move folder.", nil)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder moved successfully.", nil)
}

// UploadFileMetadataHandler godoc
//
// @Summary Upload metadata and create a file metadata in the database before upload operation.
// @Description Create a file document in the database for file upload operation.
// @Tags Files
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.UploadFileMetadataRequest true "Upload File Metadata Request"
// @Success 200 {object} models.UploadFileMetadataResponse
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /folders/{folderId}/upload [post]
func (fc *FolderController) UploadFileMetadataHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	// Define the request body structure
	var request models.UploadFileMetadataRequest

	// Bind the request body to the structure
	err := c.ShouldBindJSON(&request)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	// Cast the folder ID to ObjectID
	fileIdHex, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid folder ID.", nil)
		return
	}

	// Get information from the request context
	ownerIdHex := c.MustGet("x-user-id-hex").(primitive.ObjectID)
	ownerUsername := c.MustGet("x-username").(string)
	ownerEmail := c.MustGet("x-email").(string)

	// Get the file extension from the file name
	fileExtension := filepath.Ext(request.FileName)

	// Create the file object
	file := &models.File{
		OwnerID:        ownerIdHex,
		ParentFolderID: fileIdHex,
		FileName:       request.FileName,
		Size:           request.FileSize,
		Extension:      fileExtension,
		MimeType:       request.MimeType,
		IsDeleted:      false,
		DeletedAt:      nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	fileMetadata, err := fc.FileService.UploadFileMetadata(c, file)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to upload file metadata.", nil)
		return
	}

	// Create the response object
	response := &models.UploadFileMetadataResponse{
		File: models.FileResponse{
			ID:             fileMetadata.ID.Hex(),
			ParentFolderID: fileMetadata.ParentFolderID.Hex(),
			OwnerUsername:  ownerUsername,
			OwnerEmail:     ownerEmail,
			Name:           fileMetadata.FileName,
			MimeType:       fileMetadata.MimeType,
			Size:           fileMetadata.Size,
			CreatedAt:      fileMetadata.CreatedAt,
			UpdatedAt:      fileMetadata.UpdatedAt,
		},
		UploadURL: "",
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "File metadata uploaded successfully.", response)
}
