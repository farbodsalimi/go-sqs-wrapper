package util

import (
	"os"
)

// LoadSettings loads settings
func LoadSettings() {
	// TODO
}

// GetEnv returns the value of env variables
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
