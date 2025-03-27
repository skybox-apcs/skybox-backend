package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginFunc is the handler for the login route
func LoginFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login",
	})
}
