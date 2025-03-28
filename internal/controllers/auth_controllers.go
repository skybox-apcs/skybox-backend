package controllers

import (
	"net/http"
	"time"

	"skybox-backend/configs"
	"skybox-backend/internal/models"
	"skybox-backend/internal/services"
	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(as *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: as,
	}
}

// LoginHandler godoc
//	@Summary		Authenticates the user
//	@Description	Authenticates the user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email		body		string	true	"Email"
//	@Param			password	body		string	true	"Password"
//	@Success		200			{object}	models.User
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Invalid credentials"
//	@Router			/auth/login [post]
func (ac *AuthController) LoginHandler(c *gin.Context) {
	// Define the request body struct
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user by email
	user, err := ac.AuthService.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials."})
		return
	}

	// Compare the password with the password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials."})
		return
	}

	// Create an access token and refresh token for the user
	accessToken, err := utils.CreateAccessToken(user, configs.Config.JWTSecret, 24*7) // 7 days
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the access token."})
		return
	}

	refreshToken, err := utils.CreateRefreshToken(user, configs.Config.JWTSecret, 24*14) // 14 days
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the refresh token."})
		return
	}

	var response struct {
		User         *models.User
		AccessToken  string
		RefreshToken string
	}

	response.User = user
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// RegisterHandler is a handler that registers a new user
// RegisterHandler godoc
//	@Summary		Registers a new user
//	@Description	Registers a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email		body		string	true	"Email"
//	@Param			password	body		string	true	"Password"
//	@Param			username	body		string	true	"Username"
//	@Success		201			{string}	string	"User registered successfully"
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		409			{string}	string	"User already exists with the email"
//	@Failure		500			{string}	string	"Failed to register the user"
//	@Router			/auth/register [post]
func (ac *AuthController) RegisterHandler(c *gin.Context) {
	// Define the request body struct
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Username string `json:"username" binding:"required,max=32"`
	}

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user by email to check if the user already exists
	user, err := ac.AuthService.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists with the email."})
		return
	}

	// Encrypt the password using bcrypt
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt the password."})
		return
	}

	// Create the user object
	request.Password = string(encryptedPassword)
	user = &models.User{
		Email:        request.Email,
		PasswordHash: request.Password,
		Username:     request.Username,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Register the user
	err = ac.AuthService.RegisterUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register the user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully."})
}
