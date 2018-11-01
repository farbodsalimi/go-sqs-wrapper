package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadSettings loads settings
func LoadSettings() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv returns the value of env variables
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
