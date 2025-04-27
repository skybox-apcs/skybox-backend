package middlewares

import (
	"net/http"
	"strings"

	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

func GlobalErrorMiddleware() gin.HandlerFunc {
	// Set up a global error handler for the application
	// This middleware will catch any errors that occur during the request lifecycle
	// and return a JSON response with the error message and status code

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			errorMessage := err.Error()

			if strings.Contains(errorMessage, "not found") {
				// Handle not found errors
				shared.ErrorJSON(c, http.StatusNotFound, errorMessage)
			} else if strings.Contains(errorMessage, "permission") || strings.Contains(errorMessage, "cannot") {
				// Handle permission errors
				shared.ErrorJSON(c, http.StatusForbidden, errorMessage)
			} else if strings.Contains(errorMessage, "validation") || strings.Contains(errorMessage, "invalid") || strings.Contains(errorMessage, "required") {
				// Handle validation errors
				shared.ErrorJSON(c, http.StatusBadRequest, errorMessage)
			} else if strings.Contains(errorMessage, "unauthorized") {
				// Handle unauthorized errors
				shared.ErrorJSON(c, http.StatusUnauthorized, errorMessage)
			} else {
				// Handle system errors
				shared.ErrorJSON(c, http.StatusInternalServerError, "Internal server error. Error: "+errorMessage)
			}
		}
	}
}
