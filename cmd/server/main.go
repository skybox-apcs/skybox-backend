package main

import (
	"skybox-backend/configs"

	"skybox-backend/internal/app"
)

func main() {
	// Load Configurations
	configs.LoadConfig()

	// Initialize the application
	application, err := app.NewApp()

	if err != nil {
		panic(err)
	}

	// Run the application
	application.Run()
}
