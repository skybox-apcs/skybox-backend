package app

import (
	"fmt"
	"log"
	"net/http"

	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/routes"

	"github.com/gin-gonic/gin"
)

// Server struct encapsulates the HTTP Server
type Server struct {
	app *gin.Engine
}

// NewServer creates a new instance of the Server
func NewServer() *Server {
	// Set the release mode
	if configs.Config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		app: gin.Default(),
	}

	return s
}

// securityMiddleware sets the security headers and policies
func (s *Server) SecurityMiddleware() {
	// Set the security headers
	// X-Content-Type-Options: nosniff - Prevents browsers from MIME-sniffing a response away from the declared content-type
	// X-Frame-Options: DENY - Prevents clickjacking attacks
	// X-XSS-Protection: 1; mode=block - Prevents reflected XSS attacks
	// Strict-Transport-Security: max-age=31536000; includeSubDomains - Forces the browser to use HTTPS for the next year
	s.app.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	})

	// Set the trusted proxies (default: Disable)
	s.app.SetTrustedProxies(nil)
}

// routeMiddleware sets up the routes and the corresponding handlers
func (s *Server) RouteMiddleware() {
	routes.SetupRouter(s.app)
}

func (s *Server) GlobalErrorHandler() {
	s.app.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0]
			log.Printf("Error: %v", err)

			// Send a generic error response to the client
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
	})
}

// corsMiddleware sets up CORS headers to allow all origins
func (s *Server) CorsMiddleware() {
	s.app.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Content-Range, Range")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// StartServer starts the HTTP server
func (s *Server) StartServer() {
	host := configs.Config.BlockServerHost
	port := configs.Config.BlockServerPort

	fmt.Printf("Block Server is running on %s:%s\n", host, port)
	log.Fatal(s.app.Run(fmt.Sprintf("%s:%s", host, port)))
}
