package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		token := c.GetHeader("Authorization")
		t := strings.Split(token, " ")

		// Check if the token is valid
		if len(t) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (Missing Bearer)"})
			c.Abort()
			return
		}

		// fmt.Println("Token:", t[1])
		// fmt.Println("Secret:", secret)
		// Validate the token
		authToken := t[1]
		authorized, err := utils.IsAuthorized(authToken, secret)

		if err != nil || !authorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (Unauthorized)"})
			c.Abort()
			return
		}

		userId, err := utils.GetIDFromToken(authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (ID not found)"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("x-user-id", userId)
		c.Next()
	}
}
