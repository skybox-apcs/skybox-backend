package shared

import (
	"github.com/gin-gonic/gin"
)

// RespondJson is a helper function to send JSON responses
// It will encapsulate the response in a standard format
func RespondJson(c *gin.Context, httpStatus int, status string, message string, data any) {
	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
