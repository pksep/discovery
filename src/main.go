package main

import (
	"discovery/src/docs"
	route_handlers "discovery/src/handlers"
	"discovery/src/utils"
	"fmt"

	"discovery/src/middlewares"

	"github.com/gin-gonic/gin"
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

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:5112"

	r := gin.Default()

	secretKey, bearerToken := utils.GetBearerTokenAndSecretKey()

	if secretKey == "" {
		return
	}

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
