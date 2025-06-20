package route_handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandler_Sucess(t *testing.T) {
	routes = make(map[string]string)
	secretKey := "my_secret"

	payload := map[string]interface{}{
		"secret_key": secretKey,
		"endpoints": map[string]string{
			"test": "http://example.com/<id>",
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	RegisterHandler(w, req, secretKey)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(data), `"status":"ok"`) {
		t.Errorf("expected ok response, got %s", string(data))
	}

	if url, exists := routes["test"]; !exists || url != "http://example.com/<id>" {
		t.Errorf("route not saved correctly: %+v", routes)
	}
}

func TestRegisterHandler_WrongSecret(t *testing.T) {
	secretKey := "my_secret"
	payload := map[string]interface{}{
		"secret_key": "wrong key",
		"endpoints":  map[string]string{"test": "http://example.com/<id>"},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	RegisterHandler(w, req, secretKey)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetHandlerSuccess(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	secretKey := "my_secret"
	req := httptest.NewRequest(http.MethodGet, "/get-url?secret_key=my_secret&endpoint=test&id=42", nil)
	w := httptest.NewRecorder()

	GetURLHandler(w, req, secretKey)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(data), "http://example.com/42") {
		t.Errorf("expected url with, got %s", string(data))
	}

}

func TestGetURL_WrongSecret(t *testing.T) {
	routes = map[string]string{
		"test": "http://example.com/<id>",
	}

	secretKey := "my_secret_key"
	req := httptest.NewRequest(http.MethodGet, "/get-url?secret_key=wrong_secret_key&endpoint=test&id=42", nil)
	w := httptest.NewRecorder()

	GetURLHandler(w, req, secretKey)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetURLHandler_EndpointNotFound(t *testing.T) {
	routes = map[string]string{
		"test": "http://examaple.com/<id>",
	}

	secretKey := "my_secret"
	req := httptest.NewRequest(http.MethodGet, "/get-url?secret_key=my_secret&endpoint=unknown", nil)
	w := httptest.NewRecorder()

	GetURLHandler(w, req, secretKey)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}

}
