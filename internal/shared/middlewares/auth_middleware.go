package middlewares

import (
	"net/http"
	"strings"

	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		// Get the user ID from the token
		userId, err := utils.GetKeyFromToken("ID", authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (ID not found)"})
			c.Abort()
			return
		}

		// Get the username and email from the token
		username, err := utils.GetKeyFromToken("Username", authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (Username not found)"})
			c.Abort()
		}
		email, err := utils.GetKeyFromToken("Email", authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (Email not found)"})
			c.Abort()
			return
		}

		// Set the user ID in the context
		userIdHex, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (Hex conversion failed)"})
			c.Abort()
			return
		}

		c.Set("x-user-id", userId)
		c.Set("x-user-id-hex", userIdHex)
		c.Set("x-username", username)
		c.Set("x-email", email)

		c.Next()
	}
}
