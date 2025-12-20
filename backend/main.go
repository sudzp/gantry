package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// Workflow defines the CI/CD pipeline structure
type Workflow struct {
	Name string `yaml:"name" json:"name"`
	On   struct {
		Push struct {
			Branches []string `yaml:"branches"`
		} `yaml:"push"`
	} `yaml:"on"`
	Jobs map[string]Job `yaml:"jobs" json:"jobs"`
}

// Job represents a single job in the workflow
type Job struct {
	RunsOn string `yaml:"runs-on" json:"runs_on"`
	Steps  []Step `yaml:"steps" json:"steps"`
	Status string `json:"status"`
	Output string `json:"output"`
}

// Step represents a single step in a job
type Step struct {
	Name string `yaml:"name" json:"name"`
	Run  string `yaml:"run" json:"run"`
}

// WorkflowRun tracks execution of a workflow
type WorkflowRun struct {
	ID           string         `json:"id"`
	WorkflowName string         `json:"workflow_name"`
	Status       string         `json:"status"` // pending, running, success, failed
	Jobs         map[string]Job `json:"jobs"`
	StartedAt    time.Time      `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	mu           sync.RWMutex
}

// Server holds the application state
type Server struct {
	workflows    map[string]*Workflow
	workflowRuns map[string]*WorkflowRun
	dockerClient *client.Client
	mu           sync.RWMutex
}

func NewServer() (*Server, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Server{
		workflows:    make(map[string]*Workflow),
		workflowRuns: make(map[string]*WorkflowRun),
		dockerClient: cli,
	}, nil
}

// ParseWorkflow parses a YAML workflow file
func (s *Server) ParseWorkflow(data []byte) (*Workflow, error) {
	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}
	return &wf, nil
}

// ExecuteWorkflow runs a workflow
func (s *Server) ExecuteWorkflow(ctx context.Context, wf *Workflow) (*WorkflowRun, error) {
	runID := fmt.Sprintf("run-%d", time.Now().Unix())

	run := &WorkflowRun{
		ID:           runID,
		WorkflowName: wf.Name,
		Status:       "running",
		Jobs:         make(map[string]Job),
		StartedAt:    time.Now(),
	}

	s.mu.Lock()
	s.workflowRuns[runID] = run
	s.mu.Unlock()

	// Execute jobs
	go func() {
		defer func() {
			now := time.Now()
			run.mu.Lock()
			run.CompletedAt = &now
			run.mu.Unlock()
		}()

		allSuccess := true
		for jobName, job := range wf.Jobs {
			log.Printf("Starting job: %s", jobName)

			job.Status = "running"
			run.mu.Lock()
			run.Jobs[jobName] = job
			run.mu.Unlock()

			output, err := s.executeJob(ctx, jobName, job)

			run.mu.Lock()
			job.Output = output
			if err != nil {
				job.Status = "failed"
				allSuccess = false
				log.Printf("Job %s failed: %v", jobName, err)
			} else {
				job.Status = "success"
				log.Printf("Job %s completed successfully", jobName)
			}
			run.Jobs[jobName] = job
			run.mu.Unlock()

			if err != nil {
				break // Stop on first failure
			}
		}

		run.mu.Lock()
		if allSuccess {
			run.Status = "success"
		} else {
			run.Status = "failed"
		}
		run.mu.Unlock()
	}()

	return run, nil
}

// executeJob runs a single job using Docker
func (s *Server) executeJob(ctx context.Context, jobName string, job Job) (string, error) {
	// Default to ubuntu if no image specified
	image := "ubuntu:latest"
	if job.RunsOn == "alpine" {
		image = "alpine:latest"
	}

	// Combine all steps into a single script
	script := "#!/bin/sh\nset -e\n"
	for _, step := range job.Steps {
		script += fmt.Sprintf("echo '=== %s ==='\n", step.Name)
		script += step.Run + "\n"
	}

	// Pull image if not present
	reader, err := s.dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	io.Copy(io.Discard, reader)
	reader.Close()

	// Create container
	resp, err := s.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/sh", "-c", script},
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := s.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for completion
	statusCh, errCh := s.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return "", fmt.Errorf("container exited with status %d", status.StatusCode)
		}
	}

	// Get logs
	out, err := s.dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	defer out.Close()

	logs, _ := io.ReadAll(out)

	// Remove container
	s.dockerClient.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})

	return string(logs), nil
}

// HTTP Handlers

func (s *Server) handleUploadWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	wf, err := s.ParseWorkflow(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse workflow: %v", err), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.workflows[wf.Name] = wf
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Workflow uploaded successfully",
		"name":    wf.Name,
	})
}

func (s *Server) handleTriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	s.mu.RLock()
	wf, exists := s.workflows[name]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	run, err := s.ExecuteWorkflow(context.Background(), wf)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(run)
}

func (s *Server) handleGetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runID := vars["id"]

	s.mu.RLock()
	run, exists := s.workflowRuns[runID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Run not found", http.StatusNotFound)
		return
	}

	run.mu.RLock()
	defer run.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(run)
}

func (s *Server) handleListRuns(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	runs := make([]*WorkflowRun, 0, len(s.workflowRuns))
	for _, run := range s.workflowRuns {
		runs = append(runs, run)
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runs)
}

func (s *Server) handleListWorkflows(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	workflows := make([]*Workflow, 0, len(s.workflows))
	for _, wf := range s.workflows {
		workflows = append(workflows, wf)
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

func corsMiddleware(next http.Handler) http.Handler {
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

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/workflows", server.handleUploadWorkflow).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/workflows", server.handleListWorkflows).Methods("GET")
	r.HandleFunc("/api/workflows/{name}/trigger", server.handleTriggerWorkflow).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/runs", server.handleListRuns).Methods("GET")
	r.HandleFunc("/api/runs/{id}", server.handleGetRun).Methods("GET")

	handler := corsMiddleware(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("========================================")
	log.Println("   _____ _____ _____ _______ ______     __")
	log.Println("  / ____|  __ \\|  __ \\__   __|  _ \\ \\   / /")
	log.Println(" | |  __| |__) | |__) | | |  | |_) \\ \\_/ / ")
	log.Println(" | | |_ |  _  /|  ___/  | |  |  _ < \\   /  ")
	log.Println(" | |__| | | \\ \\| |      | |  | |_) | | |   ")
	log.Println("  \\_____|_|  \\_\\_|      |_|  |____/  |_|   ")
	log.Println("")
	log.Println("  Gantry CI/CD Platform")
	log.Println("========================================")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
