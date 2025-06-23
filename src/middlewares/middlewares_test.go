package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"discovery/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_, bearerToken := utils.GetBearerTokenAndSecretKey()

	router := gin.New()
	router.Use(AuthMiddleware(bearerToken))

	// Тестовый хук
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	t.Run("WithoutToken_ShouldReturn401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header missing or invalid")
	})

	t.Run("WithoutToken_ShouldReturn200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusOK)
		assert.Contains(t, w.Body.String(), "success")
	})

}
