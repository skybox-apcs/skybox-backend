package controllers

import (
	"net/http"

	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// UserController is the controller for the user
type UserController struct {
	UserService *services.UserService
}

// NewUserController creates a new instance of the UserController
func NewUserController(us *services.UserService) *UserController {
	return &UserController{
		UserService: us,
	}
}

// GetUserInformationHandler is a handler that returns the user information
// GetUserInformationHandler godoc
// @Summary Returns the user information
// @Description Returns the user information
// @Security		Bearer
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Router /api/v1/user [get]
func (uc *UserController) GetUserInformationHandler(c *gin.Context) {
	// Get the user ID from the request context
	userID := c.GetString("x-user-id")

	// Get the user by ID
	user, err := uc.UserService.GetUserByID(c, userID)
	if err != nil {
		shared.ErrorJSON(c, http.StatusNotFound, "User not found")
		return
	}

	shared.SuccessJSON(c, http.StatusOK, "User information retrieved successfully", user)
}

// GetUserByIDHandler godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User ID" example(1234567890abcdef12345678)
// @Success 200 {object} models.UserResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/user/{userId} [get]
func (uc *UserController) GetUserByIDHandler(c *gin.Context) {
	// Get the user ID from the URL parameters
	userID := c.Param("userId")
	if userID == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := uc.UserService.GetUserByID(c, userID)
	if err != nil {
		shared.ErrorJSON(c, http.StatusNotFound, "User not found")
		return
	}

	// Encapsulate in the UserResponse struct
	userResponse := &models.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}

	shared.SuccessJSON(c, http.StatusOK, "User information retrieved successfully", userResponse)
}

// GetUsersByIDsHandler godoc
// @Summary Get users by IDs
// @Description Get users by IDs
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserIDListRequest true "User IDs" example([{"ids": ["1234567890abcdef12345678"]}])
// @Success 200 {object} models.UserListResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/user/ids [post]
func (uc *UserController) GetUsersByIDsHandler(c *gin.Context) {
	// Get the user IDs from the request body
	var request models.UserIDListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	users, err := uc.UserService.GetUsersByIDs(c, request.IDs)
	if err != nil {
		shared.ErrorJSON(c, http.StatusNotFound, "User not found")
		return
	}

	// Encapsulate in the UserListResponse struct
	userListResponse := &models.UserListResponse{
		Users: make([]models.UserResponse, len(users)),
		Count: int64(len(users)),
	}

	for i, user := range users {
		userListResponse.Users[i] = models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		}
	}

	// Return the user information
	shared.SuccessJSON(c, http.StatusOK, "User information retrieved successfully", userListResponse)
}

// GetUserByEmailHandler godoc
// @Summary Get user by email
// @Description Get user by email
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "User email" example(email@example.com)
// @Success 200 {object} models.UserResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/user/email/{email} [get]
func (uc *UserController) GetUserByEmailHandler(c *gin.Context) {
	// Get the email from the URL parameters
	email := c.Param("email")
	if email == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := uc.UserService.GetUserByEmail(c, email)
	if err != nil {
		shared.ErrorJSON(c, http.StatusNotFound, "User not found")
		return
	}

	// Encapsulate in the UserResponse struct
	userResponse := &models.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}

	// Return the user information
	shared.SuccessJSON(c, http.StatusOK, "User information retrieved successfully", userResponse)
}

// GetUsersByEmailsHandler godoc
// @Summary Get users by emails
// @Description Get users by emails
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserEmailListRequest true "User emails" example([{"emails": ["a@example.com"]}])
// @Success 200 {object} models.UserListResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/user/emails [post]
func (uc *UserController) GetUsersByEmailsHandler(c *gin.Context) {
	// Get the emails from the request body
	var request models.UserEmailListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	users, err := uc.UserService.GetUsersByEmails(c, request.Emails)
	if err != nil {
		shared.ErrorJSON(c, http.StatusNotFound, "User not found")
		return
	}

	// Encapsulate in the UserListResponse struct
	userListResponse := &models.UserListResponse{
		Users: make([]models.UserResponse, len(users)),
		Count: int64(len(users)),
	}

	for i, user := range users {
		userListResponse.Users[i] = models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		}
	}

	// Return the user information
	shared.SuccessJSON(c, http.StatusOK, "User information retrieved successfully", userListResponse)
}
