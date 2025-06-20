package main

import (
	"discovery/src/docs"
	route_handlers "discovery/src/handlers"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Discovery Service API
// @version 1.0
// @description This is a discovery service with dynamic URL registration and retrieval.
// @host localhost:5112
// @BasePath /

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

	// Настраиваем Swagger
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:5112"

	r := gin.Default()

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Регистрация эндпоинтов
	r.POST("/register", func(c *gin.Context) {
		route_handlers.RegisterHandlerGin(c, secretKey)
	})

	// Получение URL
	r.GET("/get-url", func(c *gin.Context) {
		route_handlers.GetURLHandlerGin(c, secretKey)
	})

	fmt.Println("Discovery Service running on port 5112...")
	r.Run(":5112")
}
