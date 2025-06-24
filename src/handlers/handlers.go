package route_handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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
}

type GetEndpointDynamicRequest struct {
	Endpoint string `json:"endpoint"`
	Dynamic  string `json:"dynamic"`
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
// @Description Регистрирует переданный список endpoints
// @Tags Discovery
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Endpoints"
// @Success 200 {object} RegisterEndpointsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /register [post]
func RegisterHandlerGin(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error", "message": "Invalid JSON"})
		return
	}

	if len(request.Endpoints) == 0 || request.Endpoints == nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "error", "message": "Endpoints is empty"})
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
// @Description Возвращает URL для зарегистрированного endpoint с необязательной заменой параметров
// @Tags Discovery
// @Accept json
// @Produce json
// @Param endpoint query string true "Запрашиваемый endpoint"
// @Param dynamic_param query string false "Динамический параметр запроса. Примеры: color=red or size=large, в прямых запросах можно использовать обычный формат"
// @Success 200 {object} GetURLResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /get-url [get]
func GetURLHandlerGin(c *gin.Context) {
	endpoint := c.Query("endpoint")

	mutex.Lock()
	url, exists := routes[endpoint]
	mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"result": "error", "message": "Endpoint not found"})
		return
	}

	for key, values := range c.Request.URL.Query() {
		if key != "endpoint" && len(values) > 0 {
			log.Printf("url: %s, key: %s, value: %s", url, key, values[0])
			if strings.Contains(values[0], "=") {
				splitedDynamicParam := strings.Split(values[0], "=")

				key = splitedDynamicParam[0]
				values[0] = splitedDynamicParam[1]
			}
			re := regexp.MustCompile(fmt.Sprintf(`<%s:[^>]+>`, key))
			url = re.ReplaceAllString(url, values[0])
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": "ok", "url": url})
}
