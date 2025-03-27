package controllers

import (
	"net/http"
	"time"

	"skybox-backend/internal/models"
	"skybox-backend/internal/services"

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

func (ac *AuthController) Register(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AuthService.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists with the email."})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt the password."})
		return
	}

	request.Password = string(encryptedPassword)

	user = &models.User{
		Email:        request.Email,
		PasswordHash: request.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = ac.AuthService.RegisterUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register the user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully."})
}
