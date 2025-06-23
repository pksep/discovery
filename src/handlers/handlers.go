package route_handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	routes = make(map[string]string)
	mutex  = &sync.Mutex{}
)

type Endpoints struct {
	Endpoint string `json:"endpoint"`
	URL      string `json:"url"`
}

type RegisterRequest struct {
	Endpoints map[string]string `json:"endpoints"`
	SecretKey string            `json:"secret_key"`
}

type RegisterEndpointsResponse struct {
	Result string `json:"result"`
}

type GetURLResponse struct {
	Result string `json:"result"`
	URL    string `json:"url"`
}

type ErrorResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// RegisterHandlerGin godoc
// @Summary Register new endpoints
// @Description Registers a list of endpoints with their URLs
// @Tags Discovery
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Endpoints and Secret Key"
// @Param secret_key query string true "Secret Key"
// @Success 200 {object} RegisterEndpointsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /register [post]
func RegisterHandlerGin(c *gin.Context, secretKey string) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error", "message": "Invalid JSON"})
		return
	}

	if request.SecretKey != secretKey {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error", "message": "Wrong secret key"})
		return
	}

	mutex.Lock()
	for endpointName, url := range request.Endpoints {
		routes[endpointName] = url
	}
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}

// GetURLHandlerGin godoc
// @Summary Get URL by endpoint
// @Description Returns URL for a registered endpoint, with optional parameter substitution
// @Tags Discovery
// @Accept json
// @Produce json
// @Param secret_key query string true "Secret Key"
// @Param endpoint query string true "Endpoint name"
// @Param dynamic query string false "Optional dynamic parameter"
// @Success 200 {object} GetURLResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /get-url [get]
func GetURLHandlerGin(c *gin.Context, secretKey string) {
	secretKeyURL := c.Query("secret_key")
	endpoint := c.Query("endpoint")

	if secretKeyURL != secretKey {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error", "message": "Wrong secret key"})
		return
	}

	mutex.Lock()
	url, exists := routes[endpoint]
	mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"result": "error", "message": "Endpoint not found"})
		return
	}

	// Подстановка параметров
	for key, values := range c.Request.URL.Query() {
		if key != "secret_key" && key != "endpoint" && len(values) > 0 {
			url = strings.ReplaceAll(url, fmt.Sprintf("<%s>", key), values[0])
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": "ok", "url": url})
}
