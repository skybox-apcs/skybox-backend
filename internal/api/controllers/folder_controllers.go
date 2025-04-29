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
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {object} models.Folder
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId} [get]
func (fc *FolderController) GetFolderHandler(c *gin.Context) {
	// Get the folder ID from the request parameters
	folderId := c.Param("folderId")
	if folderId == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	folder, err := fc.FolderService.GetFolderByID(c, folderId)
	if err != nil {
		c.Error(err)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "Folder retrieved successfully.", folder)
}

// CreateFolderHandler godoc
//
// @Summary Create a new folder inside a specified parent folder
// @Description Create a new folder inside a given parent folder. If no parent folder ID is provided, the folder will be created at the root level.
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.CreateFolderRequest true "Create Folder Request"
// @Success 201 {object} models.CreateFolderResponse
// @Failure 400 {string} string "Invalid request.
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId}/create [post]
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
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid parent folder ID. Error: "+err.Error(), nil)
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
		c.Error(err)
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
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {object} models.GetFolderContentsResponse
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId}/contents [get]
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
		if folderErr != nil {
			c.Error(folderErr)
		} else {
			c.Error(fileErr)
		}
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
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Success 200 {string} string "Folder deleted successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId} [delete]
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
		c.Error(err)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder deleted successfully.", nil)
}

// RenameFolderHandler godoc
//
// @Summary Rename a specific folder
// @Description Rename an existing folder by providing its ID and the new name. Only the folder name will be updated; the folder's contents are unaffected.
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.RenameFolderRequest true "Rename Folder Request"
// @Success 200 {string} string "Folder renamed successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId}/rename [put]
// @Router /api/v1/folders/{folderId}/rename [patch]
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
		c.Error(err)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder renamed successfully.", nil)
}

// MoveFolderHandler godoc
//
// @Summary Move a folder to a new parent folder
// @Description Move a folder from its current location to another destination folder. The folder's contents will be moved along with it.
// @Security		Bearer
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.MoveFolderRequest true "Move Folder Request"
// @Success 200 {string} string "Folder moved successfully."
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId}/move [put]
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
		c.Error(err)
		return
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "Folder moved successfully.", nil)
}

// UploadFileMetadataHandler godoc
//
// @Summary Upload metadata and create a file metadata in the database before upload operation.
// @Description Create a file document in the database for file upload operation.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID" minlength(24) maxlength(24)
// @Param request body models.UploadFileMetadataRequest true "Upload File Metadata Request"
// @Success 200 {object} models.UploadFileMetadataResponse
// @Failure 400 {string} string "Invalid request."
// @Failure 404 {string} string "Folder not found."
// @Failure 500 {string} string "Internal server error."
// @Router /api/v1/folders/{folderId}/upload [post]
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

	fileMetadata, uploadURL, err := fc.FileService.UploadFileMetadata(c, file)
	if err != nil {
		c.Error(err)
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
		UploadURL: uploadURL,
	}

	// Send a success response
	shared.RespondJson(c, http.StatusOK, "success", "File metadata uploaded successfully.", response)
}

func (fc *FolderController) CheckFolderPermission(c *gin.Context, folderID string, userID string, permission string) (bool, error) {
	// Check if the user has a specific shared permission
	sharedUser, err := fc.FolderService.GetFolderSharedUser(c, folderID, userID)
	if err == nil {
		if permission == "edit" {
			return sharedUser.Permission, nil
		}
		return true, nil // View permission
	}

	// Get tge folder owner ID and check if the user is the owner
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err == nil {
		if folder.OwnerID.Hex() == userID {
			return true, nil // Owner has all permissions
		}
	}

	if permission == "edit" {
		return false, nil // No edit permission
	}

	// If no shared permission, check if the folder is public
	isPublic, err := fc.FolderService.GetFolderShareInfo(c, folderID)
	if err != nil {
		return false, err
	}

	return isPublic, nil
}

// UpdateFolderPublicStatusHandler updates the public status of a folder (public for everyone to view or restricted to only added members).
// @Summary Update folder public status of a folder (public for everyone to view or restricted to only added members)
// @Description Updates the public status of a folder by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.UpdateFolderPublicRequest true "Update Folder Public Status Request"
// @Success 200 {string} string "Folder public status updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Folder not found"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/public-status [put]
func (fc *FolderController) UpdateFolderPublicStatusHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request models.UpdateFolderPublicRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can modify this folder.", nil)
		return
	}

	err = fc.FolderService.UpdateFolderPublicStatus(c, folderID, request.IsPublic)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to update folder public status.", nil)
		return
	}

	shared.RespondJson(c, http.StatusOK, "success", "Folder public status updated successfully.", request)
}

// UpdateFolderAndSubfoldersPublicStatusHandler updates the public status of a folder and its subfolders.
// @Summary Update folder and subfolders public status
// @Description Updates the public status of a folder and all its subfolders by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.UpdateFolderPublicRequest true "Update Folder Public Status Request"
// @Success 200 {string} string "Folder and subfolders public status updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Folder not found"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/public-status/all [put]
func (fc *FolderController) UpdateFolderAndSubfoldersPublicStatusHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request models.UpdateFolderPublicRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can modify this folder.", nil)
		return
	}

	err = fc.FolderService.UpdateFolderAndSubfoldersPublicStatus(c, folderID, request.IsPublic)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to update folder and subfolders public status.", nil)
		return
	}

	shared.RespondJson(c, http.StatusOK, "success", "Folder and subfolders public status updated successfully.", request)
}

// GetFolderPublicStatusHandler retrieves the public status of a folder.
// @Summary Get folder public status
// @Description Retrieves the public status of a folder by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Success 200 {string} string "Folder public status retrieved successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Folder not found"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/public-status [get]
func (fc *FolderController) GetFolderPublicStatusHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to retrieve folder public status.", nil)
		return
	}

	shared.RespondJson(c, http.StatusOK, "success", "Folder public status retrieved successfully.", gin.H{
		"is_public": folder.IsPublic,
	})
}

// ShareFolderHandler shares a folder with a user (share edit permission by specifying permission = true, view permission by specifying permission = false.
// @Summary Share a folder
// @Description Shares a folder with a user by providing the user ID and permission.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.ShareFolderRequest true "Share Folder Request"
// @Success 200 {string} string "Folder shared successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 403 {string} string "Only the owner can share this folder"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/share [post]
func (fc *FolderController) ShareFolderHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request struct {
		UserID     string `json:"user_id"`
		Permission bool   `json:"permission"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can share this folder.", nil)
		return
	}

	err = fc.FolderService.ShareFolder(c, folderID, request.UserID, request.Permission)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to share folder.", nil)
		return
	}

	// the data return should be the list of shared users
	sharedUsers, err := fc.FolderService.GetFolderSharedUsers(c, folderID)

	shared.RespondJson(c, http.StatusOK, "success", "Folder shared successfully.", sharedUsers)
}

// RemoveFolderShareHandler removes sharing permissions for a folder.
// @Summary Remove folder share
// @Description Removes sharing permissions for a folder by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.RemoveFolderShareRequest true "Remove Folder Share Request"
// @Success 200 {string} string "Folder share removed successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 403 {string} string "Only the owner can remove share permissions"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/share [delete]
func (fc *FolderController) RemoveFolderShareHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can remove share permissions.", nil)
		return
	}
	sharedUsers, err := fc.FolderService.GetFolderSharedUsers(c, folderID)

	err = fc.FolderService.RemoveFolderShare(c, folderID, request.UserID)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to remove folder share.", nil)
		return
	}

	shared.RespondJson(c, http.StatusOK, "success", "Folder share removed successfully.", sharedUsers)
}

// GetFolderSharedUsersHandler retrieves the list of users a folder is shared with.
// @Summary Get folder shared users
// @Description Retrieves the list of users a folder is shared with by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Success 200 {array} models.FolderSharedUser "List of shared users"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Folder not found"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/shared-users [get]
func (fc *FolderController) GetFolderSharedUsersHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	sharedUsers, err := fc.FolderService.GetFolderSharedUsers(c, folderID)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to retrieve shared users.", nil)
		return
	}

	if sharedUsers == nil {
		sharedUsers = []*models.FolderSharedUser{}
	}

	shared.RespondJson(c, http.StatusOK, "success", "Shared users retrieved successfully.", sharedUsers)
}

// ShareFolderAndSubfoldersHandler shares a folder and its subfolders with a user.
// @Summary Share folder and subfolders
// @Description Shares a folder and all its subfolders with a user by providing the user ID and permission.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.ShareFolderRequest true "Share Folder Request"
// @Success 200 {string} string "Folder and subfolders shared successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 403 {string} string "Only the owner can share this folder"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/share/all [post]
func (fc *FolderController) ShareFolderAndSubfoldersHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request struct {
		UserID     string `json:"user_id"`
		Permission bool   `json:"permission"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can share this folder.", nil)
		return
	}

	err = fc.FolderService.ShareFolderAndSubfolders(c, folderID, request.UserID, request.Permission)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to share folder and subfolders.", nil)
		return
	}

	sharedUsers, err := fc.FolderService.GetFolderSharedUsers(c, folderID)
	shared.RespondJson(c, http.StatusOK, "success", "Folder and subfolders shared successfully.", sharedUsers)
}

// RevokeFolderAndSubfoldersShareHandler revokes sharing permissions for a folder and its subfolders.
// @Summary Revoke folder and subfolders share
// @Description Revokes sharing permissions for a folder and all its subfolders by its ID.
// @Tags Folders
// @Accept json
// @Produce json
// @Param folderId path string true "Folder ID"
// @Param request body models.RevokeFolderShareRequest true "Revoke Folder Share Request"
// @Success 200 {string} string "Folder and subfolders share revoked successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 403 {string} string "Only the owner can revoke share permissions"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /api/v1/folders/{folderId}/share/all [delete]
func (fc *FolderController) RevokeFolderAndSubfoldersShareHandler(c *gin.Context) {
	folderID := c.Param("folderId")
	if folderID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Folder ID is required.", nil)
		return
	}

	var request struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request.", nil)
		return
	}

	ownerID := c.MustGet("x-user-id-hex").(primitive.ObjectID).Hex()
	folder, err := fc.FolderService.GetFolderByID(c, folderID)
	if err != nil || folder.OwnerID.Hex() != ownerID {
		shared.RespondJson(c, http.StatusForbidden, "error", "Only the owner can revoke share permissions.", nil)
		return
	}
	sharedUsers, err := fc.FolderService.GetFolderSharedUsers(c, folderID)

	err = fc.FolderService.RevokeFolderAndSubfoldersShare(c, folderID, request.UserID)
	if err != nil {
		shared.RespondJson(c, http.StatusInternalServerError, "error", "Failed to revoke folder and subfolders share.", nil)
		return
	}

	shared.RespondJson(c, http.StatusOK, "success", "Folder and subfolders share revoked successfully.", sharedUsers)
}
