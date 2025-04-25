// @title           Skybox Block Server API
// @version         1.0
// @description     Skybox is a cloud-based file storage provider similar to Google Drive and Dropbox. It allows users to securely store, manage, and retrieve their files.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Skybox Support
// @contact.url    http://skybox.io/support
// @contact.email  support@skybox.io

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description  Use "Bearer {your_token}" to authenticate requests.

// @externalDocs.description  OpenAPI Documentation
// @externalDocs.url          https://swagger.io/resources/open-api/

package main

import (
	"skybox-backend/configs"
	"skybox-backend/docs"
	"skybox-backend/internal/blockserver/app"
)

func main() {
	// Load Configurations
	configs.LoadConfig()

	// Initialize Swagger
	docs.SwaggerInfo.Host = configs.Config.BlockServerHost + ":" + configs.Config.BlockServerPort

	// Start the server
	app.StartServer()
}
