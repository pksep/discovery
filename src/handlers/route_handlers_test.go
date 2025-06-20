package route_handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter(secretKey string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		RegisterHandlerGin(c, secretKey)
	})

	r.GET("/get-url", func(c *gin.Context) {
		GetURLHandlerGin(c, secretKey)
	})

	return r
}

func TestRegisterHandlerGin_Success(t *testing.T) {
	routes = make(map[string]string)
	secretKey := "my_secret"
	router := setupRouter(secretKey)

	payload := map[string]interface{}{
		"secret_key": secretKey,
		"endpoints": map[string]string{
			"test": "http://example.com/<id>",
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(data), `"result":"ok"`) && !strings.Contains(string(data), `"result": "ok"`) {
		t.Errorf("expected ok response, got %s", string(data))
	}

	if url, exists := routes["test"]; !exists || url != "http://example.com/<id>" {
		t.Errorf("route not saved correctly: %+v", routes)
	}
}

func TestRegisterHandlerGin_WrongSecret(t *testing.T) {
	secretKey := "my_secret"
	router := setupRouter(secretKey)

	payload := map[string]interface{}{
		"secret_key": "wrong_key",
		"endpoints": map[string]string{
			"test": "http://example.com/<id>",
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetURLHandlerGin_Success(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	secretKey := "my_secret"
	router := setupRouter(secretKey)

	req, _ := http.NewRequest(http.MethodGet, "/get-url?secret_key=my_secret&endpoint=test&id=42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(data), "http://example.com/42") {
		t.Errorf("expected url with parameter, got %s", string(data))
	}
}

func TestGetURLHandlerGin_WrongSecret(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	secretKey := "my_secret"
	router := setupRouter(secretKey)

	req, _ := http.NewRequest(http.MethodGet, "/get-url?secret_key=wrong_secret&endpoint=test&id=42", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetURLHandlerGin_EndpointNotFound(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	secretKey := "my_secret"
	router := setupRouter(secretKey)

	req, _ := http.NewRequest(http.MethodGet, "/get-url?secret_key=my_secret&endpoint=unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
