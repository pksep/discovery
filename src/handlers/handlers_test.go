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

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		RegisterHandlerGin(c)
	})

	r.GET("/get-url", func(c *gin.Context) {
		GetURLHandlerGin(c)
	})

	return r
}

func TestRegisterHandlerGin_Success(t *testing.T) {
	routes = make(map[string]string)
	router := setupRouter()

	payload := map[string]interface{}{
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

func TestRegisterHandlerGin_ErrorEmptyEndpoints(t *testing.T) {
	routes = make(map[string]string)
	router := setupRouter()

	payload := map[string]interface{}{}

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

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(data), `error`) {
		t.Errorf("expected error response, got %s", string(data))
	}
}

func TestGetURLHandlerGin_Success(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id: int>",
	}

	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/get-url?endpoint=test&id=42", nil)
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

func TestGetURLHandlerGin_EndpointNotFound(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/get-url?endpoint=unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
