package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	// Server Config
	ServerPort      string
	ServerHost      string
	BlockServerPort string
	BlockServerHost string
	ReleaseMode     bool

	// Allowed Origins
	AllowedOrigins []string

	// Database Config
	MongoURI    string
	MongoDBName string

	// JWT Config
	JWTSecret string

	// AWS Config
	AWSKey    string
	AWSSecret string
	AWSBucket string
	AWSRegion string
}

var Config AppConfig = AppConfig{}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file. Using default values.")
	}

	// Load the environment variables
	// Server Config
	Config.ServerPort = getEnv("SERVER_PORT", "8080")
	Config.ServerHost = getEnv("SERVER_HOST", "localhost")
	Config.BlockServerPort = getEnv("BLOCK_SERVER_PORT", "8081")
	Config.BlockServerHost = getEnv("BLOCK_SERVER_HOST", "localhost")
	Config.ReleaseMode = getEnv("GIN_MODE", "debug") == "release"

	// Allowed Origins
	Config.AllowedOrigins = strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ",")

	// Database Config
	Config.MongoURI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	Config.MongoDBName = getEnv("MONGODB_NAME", "test")

	// JWT Config
	Config.JWTSecret = getEnv("JWT_SECRET_KEY", "secret")

	// AWS Config
	Config.AWSKey = getEnv("AWS_SECRET_KEY_ID", "")
	Config.AWSSecret = getEnv("AWS_SECRET_ACCESS_KEY", "")
	Config.AWSBucket = getEnv("AWS_BUCKET_NAME", "")
	Config.AWSRegion = getEnv("AWS_REGION", "us-east-1")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
