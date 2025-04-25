package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorld is a simple handler that returns a hello world message
// HelloWorldHandler godoc
// @Summary Returns a hello world message
// @Description Returns a hello world message
// @Tags Misc
// @Accept json
// @Produce json
// @Success 200 {string} string "Hello World"
// @Router /hello [get]
func HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
