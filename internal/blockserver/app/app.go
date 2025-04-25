package app

type Application struct{}

func NewApplication() Application {
	app := &Application{}
	return *app
}

func StartServer() {
	// Create a new application
	// application := NewApplication()

	// Start the server
	ginServer := NewServer()
	ginServer.CorsMiddleware()
	ginServer.SecurityMiddleware()
	ginServer.RouteMiddleware()
	ginServer.GlobalErrorHandler()

	ginServer.StartServer()
}
