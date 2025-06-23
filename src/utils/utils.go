package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetBearerTokenAndSecretKey() (string, string) {
	execPath, _ := os.Getwd()
	projectRoot := filepath.Join(execPath)
	envPath := filepath.Join(projectRoot, "/app/.env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
		return "", ""
	}

	secretKey := os.Getenv("SECRET_KEY_DISCOVERY")
	if len(secretKey) == 0 {
		log.Printf("Error: SECRET_KEY_DISCOVERY not found")
		return "", ""
	}

	bearerToken := os.Getenv("BEARER_TOKEN")

	if len(secretKey) == 0 {
		log.Printf("Error: BEARER_TOKEN key not found")
		return "", ""
	}

	return secretKey, bearerToken
}
