package config

import (
	"log"
	"os"
	"errors"

	"github.com/joho/godotenv"
)

func Getenv(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err.Error())
		return "", err
	}
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
		return "", errors.New("Environment variable not set")
	}
	return value, nil
}