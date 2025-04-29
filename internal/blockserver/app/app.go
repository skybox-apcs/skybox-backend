package app

import (
	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/storage"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Application struct {
	s3Client *s3.Client
}

func NewApplication() Application {
	app := &Application{}

	// Connect to AWS S3
	app.s3Client = storage.GetS3Client()
	if app.s3Client == nil && configs.Config.AWSEnabled {
		panic("Failed to create AWS S3 client")
	}

	return *app
}

func (app *Application) CloseAWSClient() {
	storage.CloseAWSClient()
	// Close AWS client if needed (no explicit close method for AWS S3 client)
}

func StartServer() {
	// Create a new application
	application := NewApplication() // Uncommented and fixed the function call
	defer application.CloseAWSClient()

	// Start the server
	ginServer := NewServer()
	ginServer.CorsMiddleware()
	ginServer.SecurityMiddleware()
	ginServer.RouteMiddleware()
	ginServer.GlobalErrorHandler()

	ginServer.StartServer()
}
