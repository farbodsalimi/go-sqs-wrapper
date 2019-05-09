package util

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// GetEnvParams struct
type GetEnvParams struct {
	Key          string
	DefaultValue string
}

// ErrEnvVarEmpty returns custom error for empty or undefined environment variable
func ErrEnvVarEmpty(key string) error {
	return fmt.Errorf("getenv: %s environment variable empty", key)
}

// LoadDotEnv loads all the environment variables from the .env file
func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnvStr returns environment variable in string
func GetEnvStr(params GetEnvParams) (string, error) {
	v, exists := os.LookupEnv(params.Key)
	if exists {
		return v, nil
	} else if !exists && params.DefaultValue != "" {
		return params.DefaultValue, nil
	}
	return v, ErrEnvVarEmpty(params.Key)
}

// GetEnvInt returns environment variable in interger
func GetEnvInt(params GetEnvParams) (int, error) {
	s, err := GetEnvStr(params)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// GetEnvInt64 returns environment variable in interger 64
func GetEnvInt64(params GetEnvParams) (int64, error) {
	v, err := GetEnvInt(params)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// GetEnvBool returns environment variable in boolean
func GetEnvBool(params GetEnvParams) (bool, error) {
	s, err := GetEnvStr(params)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
