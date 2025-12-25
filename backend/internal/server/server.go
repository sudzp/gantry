package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gantry/internal/executor"
	"gantry/internal/models"
	"gantry/internal/parser"
	"gantry/internal/storage"

	"github.com/joho/godotenv"
)

// Config holds server configuration
type Config struct {
	StorageType string // "memory" or "mongodb"
	MongoURI    string
	MongoDB     string
}

// Server coordinates all components
type Server struct {
	storage  storage.Storage
	executor executor.Executor
	parser   *parser.Parser
}

// NewServer creates a new server instance
func NewServer(cfg *Config) (*Server, error) {
	// Initialize storage based on configuration
	var store storage.Storage
	var err error

	log.Println(cfg.StorageType)

	if cfg.StorageType == "mongodb" {
		log.Printf("Initializing MongoDB storage: %s/%s", cfg.MongoURI, cfg.MongoDB)
		store, err = storage.NewMongoStorage(cfg.MongoURI, cfg.MongoDB)
		if err != nil {
			return nil, fmt.Errorf("failed to create MongoDB storage: %w", err)
		}
		log.Println("✓ MongoDB connected successfully")
	} else {
		log.Println("Using in-memory storage")
		store = storage.NewMemoryStorage()
	}

	// Initialize executor
	exec, err := executor.NewDockerExecutor()
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	// Initialize parser
	p := parser.NewParser()

	return &Server{
		storage:  store,
		executor: exec,
		parser:   p,
	}, nil
}

// NewServerFromEnv creates a server from environment variables
func NewServerFromEnv() (*Server, error) {
	_ = godotenv.Load() // Loads the .env file automatically

	cfg := &Config{
		StorageType: getEnv("STORAGE_TYPE", "memory"), // "memory" or "mongodb"
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:     getEnv("MONGO_DATABASE", "gantry"),
	}

	log.Println(cfg.StorageType)

	return NewServer(cfg)
}

func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ParseAndSaveWorkflow parses and saves a workflow
func (s *Server) ParseAndSaveWorkflow(data []byte) (*models.Workflow, error) {
	wf, err := s.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	if err := s.parser.Validate(wf); err != nil {
		return nil, err
	}

	if err := s.storage.SaveWorkflow(wf); err != nil {
		return nil, err
	}

	return wf, nil
}

// ListWorkflows returns all workflows
func (s *Server) ListWorkflows() ([]*models.Workflow, error) {
	return s.storage.ListWorkflows()
}

// TriggerWorkflow triggers a workflow execution
func (s *Server) TriggerWorkflow(ctx context.Context, name string) (*models.WorkflowRun, error) {
	wf, err := s.storage.GetWorkflow(name)
	if err != nil {
		return nil, err
	}

	return s.executeWorkflow(ctx, wf)
}

// GetRun retrieves a workflow run
func (s *Server) GetRun(id string) (*models.WorkflowRun, error) {
	return s.storage.GetRun(id)
}

// ListRuns returns all workflow runs
func (s *Server) ListRuns() ([]*models.WorkflowRun, error) {
	return s.storage.ListRuns()
}

// GetWorkflowStats returns statistics for a workflow
func (s *Server) GetWorkflowStats(workflowName string) (map[string]interface{}, error) {
	runs, err := s.storage.ListRuns()
	if err != nil {
		return nil, err
	}

	// Filter runs for this workflow
	var workflowRuns []*models.WorkflowRun
	for _, run := range runs {
		if run.WorkflowName == workflowName {
			workflowRuns = append(workflowRuns, run)
		}
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_runs":       len(workflowRuns),
		"successful_runs":  0,
		"failed_runs":      0,
		"average_duration": 0,
	}

	if len(workflowRuns) == 0 {
		return stats, nil
	}

	var successCount, failureCount int
	var totalDuration int64

	for _, run := range workflowRuns {
		if run.Status == "success" {
			successCount++
		} else if run.Status == "failed" {
			failureCount++
		}

		if run.CompletedAt != nil {
			duration := run.CompletedAt.Unix() - run.StartedAt.Unix()
			totalDuration += duration
		}
	}

	stats["successful_runs"] = successCount
	stats["failed_runs"] = failureCount

	if successCount > 0 {
		stats["success_rate"] = float64(successCount) / float64(len(workflowRuns)) * 100
	}

	if len(workflowRuns) > 0 {
		stats["average_duration"] = totalDuration / int64(len(workflowRuns))
	}

	return stats, nil
}

// GetWorkflowRuns returns all runs for a specific workflow
func (s *Server) GetWorkflowRuns(workflowName string) ([]*models.WorkflowRun, error) {
	runs, err := s.storage.ListRuns()
	if err != nil {
		return nil, err
	}

	var workflowRuns []*models.WorkflowRun
	for _, run := range runs {
		if run.WorkflowName == workflowName {
			workflowRuns = append(workflowRuns, run)
		}
	}

	return workflowRuns, nil
}

// DeleteWorkflow deletes a workflow and all associated runs
func (s *Server) DeleteWorkflow(name string) error {
	// Delete all runs for this workflow (cascade delete)
	if err := s.storage.DeleteRunsByWorkflow(name); err != nil {
		log.Printf("WARNING: failed to delete runs for workflow '%s': %v", name, err)
	}

	// Delete the workflow itself
	return s.storage.DeleteWorkflow(name)
}

// executeWorkflow executes a workflow
func (s *Server) executeWorkflow(ctx context.Context, wf *models.Workflow) (*models.WorkflowRun, error) {
	runID := fmt.Sprintf("run-%d", time.Now().Unix())

	run := &models.WorkflowRun{
		ID:           runID,
		WorkflowName: wf.Name,
		Status:       "running",
		Jobs:         make(map[string]models.Job),
		JobOrder:     wf.JobOrder,
		StartedAt:    time.Now(),
	}

	if err := s.storage.SaveRun(run); err != nil {
		return nil, err
	}

	// Execute jobs asynchronously
	go s.runJobs(ctx, run, wf)

	return run, nil
}

// runJobs executes all jobs in a workflow
func (s *Server) runJobs(_ context.Context, run *models.WorkflowRun, wf *models.Workflow) {
	defer func() {
		run.Complete()

		if err := s.storage.UpdateRun(run); err != nil {
			log.Printf("ERROR: failed to update run status in storage: %v", err)
		}
	}()

	// Create a new background context with longer timeout for job execution
	// Don't use the HTTP request context as it may timeout
	jobCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	allSuccess := true
	jobOrder := wf.JobOrder
	if len(jobOrder) == 0 {
		for name := range wf.Jobs {
			jobOrder = append(jobOrder, name)
		}
	}

	for _, jobName := range jobOrder {
		job := wf.Jobs[jobName]
		log.Printf("Starting job: %s", jobName)

		jobStartTime := time.Now()
		job.Status = "running"
		job.StartedAt = jobStartTime
		run.UpdateJob(jobName, job)

		if err := s.storage.UpdateRun(run); err != nil {
			log.Printf("ERROR: failed to update run status in storage: %v", err)
		}

		output, err := s.executor.Execute(jobCtx, jobName, job)

		jobEndTime := time.Now()
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

		run.UpdateJob(jobName, job)
		if err := s.storage.UpdateRun(run); err != nil {
			log.Printf("ERROR: failed to update run status in storage: %v", err)
		}

		if err != nil {
			break // Stop on first failure
		}
	}

	if allSuccess {
		run.SetStatus("success")
	} else {
		run.SetStatus("failed")
	}

	if err := s.storage.UpdateRun(run); err != nil {
		log.Printf("ERROR: failed to update run status in storage: %v", err)
	}
}

// Cleanup performs cleanup operations
func (s *Server) Cleanup() error {
	return s.executor.Cleanup()
}
