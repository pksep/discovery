package route_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	routes = make(map[string]string)
	mutex  = &sync.Mutex{}
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, secertKey string) {
	w.Header().Set("Content-Type", "application/json")

	if len(secertKey) <= 0 {
		http.Error(w, `{"status:"error", "message":"seckretKey was not transmitted in Handler"}`, http.StatusBadRequest)
	}

	var request struct {
		Endpoints map[string]string
		SecretKey string `json:"secret_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"status":"error"}`, http.StatusBadRequest)
		return
	}

	if request.SecretKey != secertKey {
		http.Error(w, `{"status":"error", "message":"wrong secret key"}`, http.StatusBadRequest)
	}

	mutex.Lock()
	for endpoint, url := range request.Endpoints {
		routes[endpoint] = url
	}
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func GetURLHandler(w http.ResponseWriter, r *http.Request, secertKey string) {
	if len(secertKey) <= 0 {
		http.Error(w, `{"status:"error", "message":"seckretKey not found"}`, http.StatusBadRequest)
	}

	secertKeyURL := r.URL.Query().Get("secret_key")

	if secertKeyURL != secertKey {
		http.Error(w, `{"status":"error", "message":"wrong secret key"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	endpoint := r.URL.Query().Get("endpoint")

	mutex.Lock()
	url, exists := routes[endpoint]
	mutex.Unlock()

	if !exists {
		http.Error(w, `{"status": "error", "message":"Endpoint not found"}`, http.StatusNotFound)
		return
	}

	for key, values := range r.URL.Query() {
		if key != "endpoint" && len(values) > 0 {
			url = strings.ReplaceAll(url, fmt.Sprintf("<%s>", key), values[0])
		}
	}

	response := map[string]string{"url": url}

	json.NewEncoder(w).Encode(response)
}
