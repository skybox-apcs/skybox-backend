package controllers

import (
	"net/http"
	"time"

	"skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthService      *services.AuthService
	UserTokenService *services.UserTokenService
}

func NewAuthController(authService *services.AuthService, userTokenService *services.UserTokenService) *AuthController {
	return &AuthController{
		AuthService:      authService,
		UserTokenService: userTokenService,
	}
}

// respondJson is a helper function to send JSON responses
// It will encapsulate the response in a standard format
func respondJson(c *gin.Context, httpStatus int, status string, message string, data any) {
	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

// createUserToken is a helper function to create a user token in the database
// It will be used in login, refresh to add new record in the UserToken collection
func (ac *AuthController) createUserToken(c *gin.Context, user *models.User, token string, oldToken string) (string, error) {
	if oldToken != "" {
		err := ac.UserTokenService.DeleteUserToken(c, oldToken)
		if err != nil {
			return "", err
		}
	}

	// Create the user token object
	userToken := &models.UserToken{
		UserID:    user.ID,
		Token:     token,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
		ExpiredAt: time.Now().Add(time.Duration(24*14) * time.Hour),
		CreatedAt: time.Now(),
	}

	// Create the user token in the database
	err := ac.UserTokenService.CreateUserToken(c, userToken)
	if err != nil {
		return "", err
	}

	// Return the user token ID
	return userToken.ID.Hex(), nil
}

// LoginHandler godoc
//
//			@Summary		Authenticates the user
//	   @Description	This endpoint authenticates the user by checking the email and password. If the credentials are valid, it generates an access token and a refresh token.
//			@Tags			Authentication
//			@Accept			json
//			@Produce		json
//		  @Param			request body	models.LoginRequest	true	"Login Request"
//			@Success		200			{object}	models.LoginResponse	"User authenticated successfully"
//			@Failure		400			{string}	string	"Invalid request"
//			@Failure		401			{string}	string	"Invalid credentials"
//			@Router			/api/v1/auth/login [post]
func (ac *AuthController) LoginHandler(c *gin.Context) {
	// Define the request and response body structs
	var request models.LoginRequest

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		respondJson(c, http.StatusBadRequest, "error", "Invalid request. Check email or password field.", nil)
		return
	}

	// Get the user by email
	user, err := ac.AuthService.GetUserByEmail(c, request.Email)
	if err != nil {
		respondJson(c, http.StatusUnauthorized, "error", "Invalid credentials", nil)
		return
	}

	// Compare the password with the password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		respondJson(c, http.StatusUnauthorized, "error", "Invalid credentials", nil)
		return
	}

	// Create an access token and refresh token for the user
	accessToken, err := utils.CreateAccessToken(user, configs.Config.JWTSecret, 24*7) // 7 days
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the access token.", nil)
		return
	}

	refreshToken, err := utils.CreateRefreshToken(user, configs.Config.JWTSecret, 24*14) // 14 days
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the refresh token.", nil)
		return
	}

	// Update the last login time
	err = ac.AuthService.UpdateUserLastLogin(c, user.ID.Hex())
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to update the last login time.", nil)
		return
	}

	// Create a user token in the database
	_, err = ac.createUserToken(c, user, refreshToken, "")
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the user token.", nil)
		return
	}

	// Encapsulate the response
	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           user.ID.Hex(),
		Username:     user.Username,
		Email:        user.Email,
		RootFolderID: user.RootFolderID.Hex(),
	}

	// Send the response
	respondJson(c, http.StatusOK, "success", "User authenticated successfully.", response)
}

// RegisterHandler is a handler that registers a new user
// RegisterHandler godoc
//
//		@Summary		Registers a new user
//		@Description	This endpoint registers a new user by creating a new user record in the database.
//		@Tags			Authentication
//		@Accept			json
//		@Produce		json
//	  @Param      request body      models.RegisterRequest true "Register Request"
//		@Success		201			{string}	string	"User registered successfully"
//		@Failure		400			{string}	string	"Invalid request"
//		@Failure		409			{string}	string	"User already exists with the email"
//		@Failure		500			{string}	string	"Failed to register the user"
//		@Router			/api/v1/auth/register [post]
func (ac *AuthController) RegisterHandler(c *gin.Context) {
	// Define the request body struct
	var request models.RegisterRequest

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		respondJson(c, http.StatusBadRequest, "error", "Invalid request. Check email, password, or username field.", nil)
		return
	}

	// Get the user by email to check if the user already exists
	_, err = ac.AuthService.GetUserByEmail(c, request.Email)
	if err == nil {
		respondJson(c, http.StatusBadRequest, "error", "User already exists with the email.", nil)
		return
	}

	// Get the user by username to check if the user already exists
	_, err = ac.AuthService.GetUserByUsername(c, request.Username)
	if err == nil {
		respondJson(c, http.StatusBadRequest, "error", "User already exists with the username.", nil)
		return
	}

	// Encrypt the password using bcrypt
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to encrypt the password.", nil)
		return
	}

	// Create the user object
	request.Password = string(encryptedPassword)
	user := &models.User{
		Email:        request.Email,
		PasswordHash: request.Password,
		Username:     request.Username,
		RootFolderID: primitive.NilObjectID, // Set to nil for now, will be set later
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Register the user
	err = ac.AuthService.RegisterUser(c, user)
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to register the user.", nil)
		return
	}

	// Send the response
	respondJson(c, http.StatusCreated, "success", "User registered successfully.", nil)
}

// RefreshHandler godoc
//
//		@Summary		Validate and refresh the access token via refresh token
//	 @Description	This endpoint validates the refresh token and generates a new access token, allowing the user to continue their session without re-authenticating.
//		@Tags			Authentication
//		@Accept			json
//		@Produce		json
//		@Param			request body	models.RefreshRequest	true	"Refresh Request"
//		@Success		200			{object}	models.RefreshResponse	"Access token refreshed successfully"
//		@Failure		400			{string}	string	"Invalid request"
//		@Failure		401			{string}	string	"Invalid refresh token"
//		@Failure		500			{string}	string	"Failed to refresh the access token"
//		@Router			/api/v1/auth/refresh [post]
func (ac *AuthController) RefreshHandler(c *gin.Context) {
	// Define the request and response body structs
	var request models.RefreshRequest

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		respondJson(c, http.StatusBadRequest, "error", "Invalid request. Check refresh token field.", nil)
		return
	}

	// Validate the refresh token
	user_id, err := utils.GetKeyFromToken("ID", request.RefreshToken, configs.Config.JWTSecret)
	if err != nil {
		respondJson(c, http.StatusUnauthorized, "error", "Invalid refresh token.", nil)
		return
	}

	// Get the user by ID
	user, err := ac.AuthService.GetUserByID(c, user_id)
	if err != nil {
		respondJson(c, http.StatusUnauthorized, "error", "Invalid refresh token.", nil)
		return
	}

	// Create a new access token and refresh token for the user
	accessToken, err := utils.CreateAccessToken(user, configs.Config.JWTSecret, 24*7) // 7 days
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the access token.", nil)
		return
	}
	refreshToken, err := utils.CreateRefreshToken(user, configs.Config.JWTSecret, 24*14) // 14 days
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the refresh token.", nil)
		return
	}

	// Update the last login time
	err = ac.AuthService.UpdateUserLastLogin(c, user.ID.Hex())
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to update the last login time.", nil)
		return
	}

	// Create a user token in the database
	_, err = ac.createUserToken(c, user, refreshToken, request.RefreshToken)
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to create the user token.", nil)
		return
	}

	// Encapsulate the response
	var response models.RefreshResponse
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken
	response.ID = user.ID.Hex()
	response.Username = user.Username
	response.Email = user.Email

	// Send the response
	respondJson(c, http.StatusOK, "success", "Access token refreshed successfully.", response)
}

// LogoutHandler godoc
//
//	@Summary		Logs out the user
//	@Description	Logs out the user and invalidates the refresh token
//	@Security		Bearer
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request body	models.LogoutRequest	true	"Logout Request"
//	@Success		200			{string}	string	"User logged out successfully"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Failed to log out the user"
//	@Router			/api/v1/auth/logout [post]
func (ac *AuthController) LogoutHandler(c *gin.Context) {
	// Define the request and response body structs
	var request models.LogoutRequest

	// Bind the request body to the struct and check if JSON object is valid
	err := c.ShouldBind(&request)
	if err != nil {
		respondJson(c, http.StatusBadRequest, "error", "Invalid request. Check refresh token field.", nil)
		return
	}

	// Validate the refresh token
	_, err = utils.GetKeyFromToken("ID", request.RefreshToken, configs.Config.JWTSecret)
	if err != nil {
		respondJson(c, http.StatusUnauthorized, "error", "Invalid refresh token.", nil)
		return
	}

	// Delete the user token from the database
	err = ac.UserTokenService.DeleteUserToken(c, request.RefreshToken)
	if err != nil {
		respondJson(c, http.StatusInternalServerError, "error", "Failed to log out the user.", nil)
		return
	}

	// Send the response
	respondJson(c, http.StatusOK, "success", "User logged out successfully.", nil)
}
