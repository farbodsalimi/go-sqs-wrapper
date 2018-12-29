package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type GetEnvParams struct {
	Key          string
	DefaultValue string
}

func ErrEnvVarEmpty(key string) error {
	return errors.New(fmt.Sprintf("getenv: %s environment variable empty", key))
}

// LoadDotEnv loads all the environment variables from the .env file
func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

//
func GetEnvStr(params GetEnvParams) (string, error) {
	v, exists := os.LookupEnv(params.Key)
	if exists {
		return v, nil
	} else if !exists && params.DefaultValue != "" {
		return params.DefaultValue, nil
	}
	return v, ErrEnvVarEmpty(params.Key)
}

//
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

func GetEnvInt64(params GetEnvParams) (int64, error) {
	v, err := GetEnvInt(params)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

//
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
