package app

import (
	"log"

	"skybox-backend/configs"
	"skybox-backend/internal/repositories"
)

type App struct {
	server *Server
}

func NewApp() (*App, error) {
	// Initalize DB Server
	db, err := repositories.NewMongoClient(configs.Config.MongoURI, "your-db-name")
	if err != nil {
		return nil, err
	}

	// Wire up the dependencies
	initDependencies(db)

	// Initialize the HTTP server
	server := NewServer()

	return &App{
		server: server,
	}, nil
}

func (a *App) Run() {
	// Start the HTTP server
	log.Println("Starting the server...")
	a.server.startServer()
}
