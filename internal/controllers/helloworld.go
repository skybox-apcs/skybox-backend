package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorld is a simple handler that returns a hello world message
func HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
