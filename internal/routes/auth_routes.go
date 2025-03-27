package routes

import (
	"skybox-backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

// Example route for the auth group
func RegisterLoginRoutes(router *gin.RouterGroup) {
	router.POST("/login", controllers.LoginFunc)
}
