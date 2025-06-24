package main

import (
	"discovery/docs"
	route_handlers "discovery/handlers"
	"discovery/utils"
	"fmt"

	"discovery/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Discovery Service API
// @version 1.0
// @description Discovery сервис, который принимает статические/динамические эндпоинты и может их возвращать.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
func main() {

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "discovery.pksep.ru"

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Можно указать конкретные домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	bearerToken := utils.GetBearerToken()

	if bearerToken == "" {
		return
	}

	authorized := r.Group("/", middlewares.AuthMiddleware(bearerToken))
	{
		authorized.POST("/register", func(c *gin.Context) {
			route_handlers.RegisterHandlerGin(c)
		})
		authorized.GET("/get-url", func(c *gin.Context) {
			route_handlers.GetURLHandlerGin(c)
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Discovery Service running on port 5112...")
	r.Run(":5112")
}
