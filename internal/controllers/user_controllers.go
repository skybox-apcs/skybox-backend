package controllers

import (
	"net/http"

	"skybox-backend/internal/services"

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
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Router /user [get]
func (uc *UserController) GetUserInformationHandler(c *gin.Context) {
	// Get the user ID from the request context
	userID := c.GetString("x-user-id")

	// Get the user by ID
	user, err := uc.UserService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}

	c.JSON(http.StatusOK, user)
}
