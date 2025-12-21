package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(h *Handler) http.Handler {
	r := mux.NewRouter()

	// Workflow routes
	r.HandleFunc("/api/workflows", h.HandleUploadWorkflow).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/workflows", h.HandleListWorkflows).Methods("GET")
	r.HandleFunc("/api/workflows/{name}/trigger", h.HandleTriggerWorkflow).Methods("POST", "OPTIONS")

	// Run routes
	r.HandleFunc("/api/runs", h.HandleListRuns).Methods("GET")
	r.HandleFunc("/api/runs/{id}", h.HandleGetRun).Methods("GET")

	// Apply middleware
	return CORSMiddleware(r)
}

// CORSMiddleware handles CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
