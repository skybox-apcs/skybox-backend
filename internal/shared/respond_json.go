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

// SuccessJSON sends a JSON response with a success message and data
func SuccessJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// ErrorJSON sends a JSON response with an error message
func ErrorJSON(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}
