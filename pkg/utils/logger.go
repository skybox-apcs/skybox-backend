package utils

import (
	"log"
	"os"
)

// InitLogger initializes a simple logger
func InitLogger() *log.Logger {
	return log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
