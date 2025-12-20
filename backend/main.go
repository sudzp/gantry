package main

import (
	"bytes"
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
	Jobs     map[string]Job `yaml:"jobs" json:"jobs"`
	JobOrder []string       `json:"job_order"` // Preserve YAML order
}

// Job represents a single job in the workflow
type Job struct {
	RunsOn    string     `yaml:"runs-on" json:"runs_on"`
	Steps     []Step     `yaml:"steps" json:"steps"`
	Status    string     `json:"status"`
	Output    string     `json:"output"`
	StartedAt time.Time  `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// Step represents a single step in a job
type Step struct {
	Name      string     `yaml:"name" json:"name"`
	Run       string     `yaml:"run" json:"run"`
	Status    string     `json:"status,omitempty"`
	StartedAt time.Time  `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Output    string     `json:"output,omitempty"`
}

// WorkflowRun tracks execution of a workflow
type WorkflowRun struct {
	ID           string         `json:"id"`
	WorkflowName string         `json:"workflow_name"`
	Status       string         `json:"status"` // pending, running, success, failed
	Jobs         map[string]Job `json:"jobs"`
	JobOrder     []string       `json:"job_order"` // Preserve execution order
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
	// First, parse the raw YAML to preserve key order
	var rawMap map[string]interface{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Now parse into struct
	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Extract job order from raw YAML (preserves insertion order)
	if jobsMap, ok := rawMap["jobs"].(map[string]interface{}); ok {
		wf.JobOrder = make([]string, 0, len(jobsMap))
		// In Go 1.12+, map iteration order is randomized but we can use
		// yaml.v3 which preserves order
		var orderedYAML struct {
			Jobs yaml.Node `yaml:"jobs"`
		}
		if err := yaml.Unmarshal(data, &orderedYAML); err == nil {
			// Extract keys in order from yaml.Node
			if orderedYAML.Jobs.Kind == yaml.MappingNode {
				for i := 0; i < len(orderedYAML.Jobs.Content); i += 2 {
					if orderedYAML.Jobs.Content[i].Value != "" {
						wf.JobOrder = append(wf.JobOrder, orderedYAML.Jobs.Content[i].Value)
					}
				}
			}
		}
	}

	// Fallback: if JobOrder is still empty, use map keys (will be random)
	if len(wf.JobOrder) == 0 {
		for jobName := range wf.Jobs {
			wf.JobOrder = append(wf.JobOrder, jobName)
		}
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
		JobOrder:     wf.JobOrder, // Preserve YAML order
		StartedAt:    time.Now(),
	}

	s.mu.Lock()
	s.workflowRuns[runID] = run
	s.mu.Unlock()

	// Execute jobs in order
	go func() {
		defer func() {
			now := time.Now()
			run.mu.Lock()
			run.CompletedAt = &now
			run.mu.Unlock()
		}()

		allSuccess := true
		// Use preserved job order from YAML
		jobOrder := wf.JobOrder
		if len(jobOrder) == 0 {
			// Fallback to map keys if order not preserved
			for name := range wf.Jobs {
				jobOrder = append(jobOrder, name)
			}
		}

		// Execute jobs sequentially in YAML order
		for _, jobName := range jobOrder {
			job := wf.Jobs[jobName]
			log.Printf("Starting job: %s", jobName)

			jobStartTime := time.Now()
			job.Status = "running"
			job.StartedAt = jobStartTime
			run.mu.Lock()
			run.Jobs[jobName] = job
			run.mu.Unlock()

			output, err := s.executeJob(ctx, jobName, job)

			jobEndTime := time.Now()
			run.mu.Lock()
			job.Output = output
			job.EndedAt = &jobEndTime
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

	// Build script with step tracking
	script := "#!/bin/sh\nset -e\n"
	for i, step := range job.Steps {
		script += fmt.Sprintf("\n# Step %d: %s\n", i+1, step.Name)
		script += fmt.Sprintf("echo '=== [' $(date '+%%Y-%%m-%%d %%H:%%M:%%S') '] Starting: %s ==='\n", step.Name)
		script += step.Run + "\n"
		script += fmt.Sprintf("echo '=== [' $(date '+%%Y-%%m-%%d %%H:%%M:%%S') '] Completed: %s ==='\n", step.Name)
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

	const banner = `
  ____             _                   
 / ___| __ _ _ __ | |_ _ __ _   _ 
| |  _ / _' | '_ \| __| '__| | | |
| |_| | (_| | | | | |_| |  | |_| |
 \____|\__,_|_| |_|\__|_|   \__, |
                            |___/ `

	log.Println("========================================")

	log.Println("========================================")
	log.Println(banner)
	log.Println("")
	log.Println("  Gantry CI/CD Platform")
	log.Println("========================================")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
