package main

import (
	"discovery/src/docs"
	route_handlers "discovery/src/handlers"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"discovery/src/middlewares"

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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
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

	bearerToken := os.Getenv("BEARER_TOKEN")
	fmt.Println(bearerToken)

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:5112"

	r := gin.Default()

	authorized := r.Group("/", middlewares.AuthMiddleware(bearerToken))
	{
		authorized.POST("/register", func(c *gin.Context) {
			route_handlers.RegisterHandlerGin(c, secretKey)
		})
		authorized.GET("/get-url", func(c *gin.Context) {
			route_handlers.GetURLHandlerGin(c, secretKey)
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Discovery Service running on port 5112...")
	r.Run(":5112")
}
