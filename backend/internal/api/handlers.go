// Package api provides HTTP handlers for the Gantry API
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gantry/internal/server"

	"github.com/gorilla/mux"
)

// Handler manages HTTP requests
type Handler struct {
	server *server.Server
}

// NewHandler creates a new API handler
func NewHandler(srv *server.Server) *Handler {
	return &Handler{
		server: srv,
	}
}

// HandleUploadWorkflow handles workflow upload requests
func (h *Handler) HandleUploadWorkflow(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	wf, err := h.server.ParseAndSaveWorkflow(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse workflow: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Workflow uploaded successfully",
		"name":    wf.Name,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleListWorkflows handles listing workflows
func (h *Handler) HandleListWorkflows(w http.ResponseWriter, _ *http.Request) {
	workflows, err := h.server.ListWorkflows()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list workflows: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(workflows); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleTriggerWorkflow handles workflow trigger requests
func (h *Handler) HandleTriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	run, err := h.server.TriggerWorkflow(r.Context(), name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to trigger workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(run); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleGetRun handles get run details requests
func (h *Handler) HandleGetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runID := vars["id"]

	run, err := h.server.GetRun(runID)
	if err != nil {
		http.Error(w, "Run not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(run); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleListRuns handles listing all runs
func (h *Handler) HandleListRuns(w http.ResponseWriter, _ *http.Request) {
	runs, err := h.server.ListRuns()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(runs); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleDeleteWorkflow handles workflow deletion
func (h *Handler) HandleDeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := h.server.DeleteWorkflow(name); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete workflow: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Workflow deleted successfully",
		"name":    name,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleGetWorkflowStats handles getting workflow statistics
func (h *Handler) HandleGetWorkflowStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	stats, err := h.server.GetWorkflowStats(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleGetWorkflowRuns handles getting workflow run history
func (h *Handler) HandleGetWorkflowRuns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	runs, err := h.server.GetWorkflowRuns(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get runs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(runs); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
