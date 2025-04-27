package app

import (
	"fmt"
	"net/http"
	"time"

	"skybox-backend/configs"
	"skybox-backend/internal/api/routes"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.mongodb.org/mongo-driver/mongo"
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

		// Additional security headers
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		c.Next()
	})

	// Set the trusted proxies (default: Disable)
	s.app.SetTrustedProxies(nil)
}

// rateLimitMiddleware implements rate limiting
func (s *Server) RateLimitMiddleware() {
	// Create a rate limiter: 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	s.app.Use(func(c *gin.Context) {
		context := c.Request.Context()
		limiterCtx, err := instance.Get(context, c.ClientIP())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limiterCtx.Reached {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	})
}

// routeMiddleware sets up the routes and the corresponding handlers
func (s *Server) RouteMiddleware(db *mongo.Database) {
	s.app = routes.SetupRouter(db, s.app)
}

// globalErrorHandler set up a centralized error handler with secure defaults
func (s *Server) GlobalErrorHandler() {
	s.app.Use(middlewares.GlobalErrorMiddleware())
}

// corsMiddleware sets up CORS headers with more secure defaults
func (s *Server) CorsMiddleware() {
	s.app.Use(func(c *gin.Context) {
		// Replace * with specific allowed origins in production
		// allowedOrigins := configs.Config.AllowedOrigins
		// if len(allowedOrigins) == 0 {
		// 	allowedOrigins = []string{"*"}
		// }

		// origin := c.GetHeader("Origin")
		// if origin == "" {
		// 	origin = "*"
		// }

		// // Check if the origin is allowed
		// isAllowed := false
		// for _, allowedOrigin := range allowedOrigins {
		// 	if allowedOrigin == "*" || allowedOrigin == origin {
		// 		isAllowed = true
		// 		break
		// 	}
		// }

		// if isAllowed {
		// 	c.Header("Access-Control-Allow-Origin", origin)
		// }

		c.Header("Access-Control-Allow-Origin", "*") // Replace * with specific allowed origins in production
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// Start initializes the server and starts listening on the specified port
func (s *Server) StartServer() {
	host := configs.Config.ServerHost
	port := configs.Config.ServerPort

	// Start the server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      s.app,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Server failed to start", err)
	}
}
