package main

import (
	route_handlers "discovery/src/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	execPath, _ := os.Getwd()
	projectRoot := filepath.Join(execPath)
	envPath := filepath.Join(projectRoot, ".env")

	err := godotenv.Load(envPath)

	if err != nil {
		log.Printf("Error loading .env file: %s", err)
		return
	}

	secretKey := os.Getenv("SECRET_KEY_DISCOVERY")

	if len(secretKey) == 0 {
		log.Printf("Error: secret key not found")
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		route_handlers.RegisterHandler(w, r, secretKey)
	})
	http.HandleFunc("/get-url", func(w http.ResponseWriter, r *http.Request) {
		route_handlers.GetURLHandler(w, r, secretKey)
	})

	fmt.Println("Discovery Service running on port 5112...")
	http.ListenAndServe(":5112", nil)
}
