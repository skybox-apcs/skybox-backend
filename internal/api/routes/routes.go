package routes

import (
	"skybox-backend/configs"
	"skybox-backend/internal/api/controllers"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/repositories"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationContainer struct {
	// Repositories
	ChunkRepository         *repositories.ChunkRepository
	FileRepository          *repositories.FileRepository
	FolderRepository        *repositories.FolderRepository
	UserRepository          *repositories.UserRepository
	UserTokenRepository     *repositories.UserTokenRepository
	UploadSessionRepository *repositories.UploadSessionRepository

	// Services
	AuthService          *services.AuthService
	ChunkService         *services.ChunkService
	FileService          *services.FileService
	FolderService        *services.FolderService
	UserService          *services.UserService
	UserTokenService     *services.UserTokenService
	UploadSessionService *services.UploadSessionService

	// Controllers
	AuthController          *controllers.AuthController
	FileController          *controllers.FileController
	FolderController        *controllers.FolderController
	UploadSessionController *controllers.UploadSessionController
	UserController          *controllers.UserController
}

func (app *ApplicationContainer) SetupRepositories(db *mongo.Database) {
	app.ChunkRepository = repositories.NewChunkRepository(db, models.CollectionChunks)
	app.FileRepository = repositories.NewFileRepository(db, models.CollectionFiles)
	app.FolderRepository = repositories.NewFolderRepository(db, models.CollectionFolders)
	app.UserRepository = repositories.NewUserRepository(db, models.CollectionUsers)
	app.UserTokenRepository = repositories.NewUserTokenRepository(db, models.CollectionUserTokens)
	app.UploadSessionRepository = repositories.NewUploadSessionRepository(db, models.CollectionUploadSessions)
}

func (app *ApplicationContainer) SetupServices() {
	app.AuthService = services.NewAuthService(app.UserRepository)
	app.ChunkService = services.NewChunkService(app.ChunkRepository)
	app.FileService = services.NewFileService(app.FileRepository, app.UploadSessionRepository)
	app.FolderService = services.NewFolderService(app.FolderRepository)
	app.UserService = services.NewUserService(app.UserRepository)
	app.UserTokenService = services.NewUserTokenService(app.UserTokenRepository)
	app.UploadSessionService = services.NewUploadSessionService(app.UploadSessionRepository)
}

func (app *ApplicationContainer) SetupControllers() {
	app.AuthController = controllers.NewAuthController(app.AuthService, app.UserTokenService)
	app.FileController = controllers.NewFileController(app.FileService)
	app.FolderController = controllers.NewFolderController(app.FolderService, app.FileService)
	app.UploadSessionController = controllers.NewUploadSessionController(app.UploadSessionService)
	app.UserController = controllers.NewUserController(app.UserService)
}

var appContainer *ApplicationContainer

// NewApplicationContainer initializes the application container with repositories, services, and controllers
func GetApplicationContainer(db *mongo.Database) *ApplicationContainer {
	// Initialize the application container only once
	if appContainer != nil {
		return appContainer
	}

	appContainer := &ApplicationContainer{}
	appContainer.SetupRepositories(db)
	appContainer.SetupServices()
	appContainer.SetupControllers()
	return appContainer
}

// SetupRoutes sets up the routes and the corresponding handlers
func SetupRouter(db *mongo.Database, gin *gin.Engine) *gin.Engine {
	// Swagger routes
	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRouter := gin.Group("")

	// Setup the v1 routes
	v1 := publicRouter.Group("/api/v1")

	// Public routes
	{
		// Setup the auth routes
		NewAuthRouters(db, v1)

		// Setup the user routes
		NewUserRouters(db, v1)

		// Hello World routes
		v1.GET("/hello", controllers.HelloWorldHandler)
	}

	// Private routes
	protectedRouter := gin.Group("")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(configs.Config.JWTSecret))

	v1 = protectedRouter.Group("/api/v1")

	{
		// Setup the folder routes
		NewFolderRouters(db, v1)

		// Setup the file routes
		NewFileRouters(db, v1)

		// Setup the upload routes
		NewUploadRouters(db, v1)

		// Setup the search routes
		NewSearchRouters(db, v1)
	}

	return gin
}
