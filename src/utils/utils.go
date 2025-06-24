package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetBearerToken() string {
	execPath, _ := os.Getwd()
	projectRoot := filepath.Join(execPath)

	start_mode := os.Getenv("START_MODE")

	var envPath string
	if start_mode == "docker" {
		envPath = filepath.Join(projectRoot, "/app/.env")
	} else {
		envPath = filepath.Join(projectRoot, ".env")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
		return ""
	}

	bearerToken := os.Getenv("BEARER_TOKEN")

	if len(bearerToken) == 0 {
		log.Printf("Error: BEARER_TOKEN key not found")
		return ""
	}

	return bearerToken
}
