package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ReleaseMode bool

	// API Server Config
	ServerPort string
	ServerHost string

	// BlockServer Config
	BlockServerPort  string
	BlockServerHost  string
	MaxWorkers       int
	DefaultChunkSize int64
	MaxChunkSize     int64

	// Allowed Origins
	AllowedOrigins []string

	// Database Config
	MongoURI    string
	MongoDBName string

	// JWT Config
	JWTSecret string

	// AWS Config
	AWSEnabled      bool
	AWSKey          string
	AWSSecret       string
	AWSSessionToken string
	AWSBucket       string
	AWSRegion       string
}

// Config is the global application configuration
// It is initialized with default values and can be overridden by environment variables
var Config AppConfig = AppConfig{
	ReleaseMode: false,

	MaxWorkers:       5,
	DefaultChunkSize: 5242880,   // 5MB
	MaxChunkSize:     104857600, // 100MB

	JWTSecret: "secret",
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file. Using default values.")
	}

	// Load the environment variables
	Config.ReleaseMode = getEnv("GIN_MODE", "debug") == "release"
	// API Server Config
	configAPIServer()

	// BlockServer Config
	configBlockServer()

	// Allowed Origins
	Config.AllowedOrigins = strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ",")

	// Database Config
	Config.MongoURI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	Config.MongoDBName = getEnv("MONGODB_NAME", "test")

	// JWT Config
	Config.JWTSecret = getEnv("JWT_SECRET_KEY", "secret")

	// AWS Config
	configAWS()

	// Print Configs
	fmt.Printf("Loaded Config: %+v\n", Config)
}

func configAPIServer() {
	Config.ServerPort = getEnv("SERVER_PORT", "8080")
	Config.ServerHost = getEnv("SERVER_HOST", "localhost")
}

func configBlockServer() {
	var err error

	Config.BlockServerPort = getEnv("BLOCK_SERVER_PORT", "8081")
	Config.BlockServerHost = getEnv("BLOCK_SERVER_HOST", "localhost")
	Config.MaxWorkers, err = strconv.Atoi(getEnv("MAX_WORKERS", "5"))
	if err != nil {
		log.Println("Invalid MAX_WORKERS value, using default value of 5")
		Config.MaxWorkers = 5
	}
	Config.DefaultChunkSize, err = strconv.ParseInt(getEnv("DEFAULT_CHUNK_SIZE", "5242880"), 10, 64) // 5MB
	if err != nil {
		log.Println("Invalid DEFAULT_CHUNK_SIZE value, using default value of 5MB")
		Config.DefaultChunkSize = 5242880 // 5MB
	}
	Config.MaxChunkSize, err = strconv.ParseInt(getEnv("MAX_CHUNK_SIZE", "104857600"), 10, 64) // 100MB
	if err != nil {
		log.Println("Invalid MAX_CHUNK_SIZE value, using default value of 100MB")
		Config.MaxChunkSize = 104857600 // 100MB
	}
	if Config.MaxChunkSize < Config.DefaultChunkSize {
		log.Println("MAX_CHUNK_SIZE should be greater than or equal to DEFAULT_CHUNK_SIZE. Setting MAX_CHUNK_SIZE to DEFAULT_CHUNK_SIZE.")
		Config.MaxChunkSize = Config.DefaultChunkSize
	}
}

func configAWS() {
	Config.AWSEnabled = getEnv("AWS_ENABLED", "false") == "true"
	Config.AWSKey = getEnv("AWS_ACCESS_KEY_ID", "")
	Config.AWSSecret = getEnv("AWS_SECRET_ACCESS_KEY", "")
	Config.AWSSessionToken = getEnv("AWS_SESSION_TOKEN", "")
	Config.AWSBucket = getEnv("AWS_S3_BUCKET_NAME", "")
	Config.AWSRegion = getEnv("AWS_REGION", "us-east-1")
}

// getEnv retrieves the value of an environment variable or returns a fallback value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
