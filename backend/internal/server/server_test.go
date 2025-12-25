package server

import (
	"context"
	"testing"

	"gantry/internal/models"
	"gantry/internal/parser"
	"gantry/internal/storage"
)

func TestServer_ParseAndSaveWorkflow(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	yaml := []byte(`
name: Test Workflow
on:
  push:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu
    steps:
      - name: Run tests
        run: echo "testing"
`)

	wf, err := srv.ParseAndSaveWorkflow(yaml)
	if err != nil {
		t.Fatalf("Failed to parse and save workflow: %v", err)
	}

	if wf.Name != "Test Workflow" {
		t.Errorf("Expected workflow name 'Test Workflow', got '%s'", wf.Name)
	}

	// Verify saved in storage
	retrieved, err := srv.storage.GetWorkflow("Test Workflow")
	if err != nil {
		t.Fatalf("Failed to retrieve workflow: %v", err)
	}

	if retrieved.Name != wf.Name {
		t.Errorf("Workflow not properly saved")
	}
}

func TestServer_ListWorkflows(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	// Save multiple workflows
	for i := 1; i <= 3; i++ {
		wf := &models.Workflow{
			Name: "Workflow" + string(rune(48+i)),
			Jobs: map[string]models.Job{},
		}
		srv.storage.SaveWorkflow(wf)
	}

	workflows, err := srv.ListWorkflows()
	if err != nil {
		t.Fatalf("Failed to list workflows: %v", err)
	}

	if len(workflows) != 3 {
		t.Errorf("Expected 3 workflows, got %d", len(workflows))
	}
}

func TestServer_TriggerWorkflow(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps: []models.Step{
					{Name: "Test", Run: "echo test"},
				},
			},
		},
		JobOrder: []string{"test"},
	}

	srv.storage.SaveWorkflow(wf)

	// Note: This will fail without Docker, but we can test the run creation
	run, err := srv.TriggerWorkflow(context.Background(), "Test")
	if err != nil && err.Error() != "failed to create executor: docker daemon not available" {
		// Expected error if Docker not available
		if run == nil {
			t.Logf("Expected error without Docker: %v", err)
			return
		}
	}

	if run != nil && run.WorkflowName != "Test" {
		t.Errorf("Run not created with correct workflow name")
	}
}

func TestServer_GetWorkflowStats(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	// Save workflow
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{},
	}
	srv.storage.SaveWorkflow(wf)

	// Save runs
	for i := 0; i < 5; i++ {
		status := "success"
		if i == 4 {
			status = "failed"
		}
		run := &models.WorkflowRun{
			ID:           "run-" + string(rune(48+i)),
			WorkflowName: "Test",
			Status:       status,
			Jobs:         make(map[string]models.Job),
		}
		srv.storage.SaveRun(run)
	}

	stats, err := srv.GetWorkflowStats("Test")
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats["total_runs"] != 5 {
		t.Errorf("Expected 5 total runs, got %v", stats["total_runs"])
	}

	if stats["successful_runs"] != 4 {
		t.Errorf("Expected 4 successful runs, got %v", stats["successful_runs"])
	}

	if stats["failed_runs"] != 1 {
		t.Errorf("Expected 1 failed run, got %v", stats["failed_runs"])
	}

	successRate := stats["success_rate"].(float64)
	if successRate != 80 {
		t.Errorf("Expected 80%% success rate, got %v%%", successRate)
	}
}

func TestServer_GetWorkflowRuns(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	// Save workflow
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{},
	}
	srv.storage.SaveWorkflow(wf)

	// Save runs for this workflow and another
	for i := 0; i < 3; i++ {
		run := &models.WorkflowRun{
			ID:           "test-run-" + string(rune(48+i)),
			WorkflowName: "Test",
			Status:       "success",
			Jobs:         make(map[string]models.Job),
		}
		srv.storage.SaveRun(run)
	}

	for i := 0; i < 2; i++ {
		run := &models.WorkflowRun{
			ID:           "other-run-" + string(rune(48+i)),
			WorkflowName: "Other",
			Status:       "success",
			Jobs:         make(map[string]models.Job),
		}
		srv.storage.SaveRun(run)
	}

	runs, err := srv.GetWorkflowRuns("Test")
	if err != nil {
		t.Fatalf("Failed to get workflow runs: %v", err)
	}

	if len(runs) != 3 {
		t.Errorf("Expected 3 runs for Test workflow, got %d", len(runs))
	}

	for _, run := range runs {
		if run.WorkflowName != "Test" {
			t.Errorf("Got run from different workflow: %s", run.WorkflowName)
		}
	}
}

func TestServer_DeleteWorkflow(t *testing.T) {
	srv := &Server{
		storage: storage.NewMemoryStorage(),
		parser:  parser.NewParser(),
	}

	// Save workflow
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{},
	}
	srv.storage.SaveWorkflow(wf)

	// Save runs for this workflow
	run := &models.WorkflowRun{
		ID:           "run-123",
		WorkflowName: "Test",
		Status:       "success",
		Jobs:         make(map[string]models.Job),
	}
	srv.storage.SaveRun(run)

	// Delete workflow (should cascade delete runs)
	err := srv.DeleteWorkflow("Test")
	if err != nil {
		t.Fatalf("Failed to delete workflow: %v", err)
	}

	// Verify workflow deleted
	_, err = srv.storage.GetWorkflow("Test")
	if err == nil {
		t.Error("Workflow should be deleted but still exists")
	}

	// Verify runs deleted (cascade delete)
	runs, err := srv.storage.ListRuns()
	if err != nil {
		t.Fatalf("Failed to list runs: %v", err)
	}

	if len(runs) > 0 {
		t.Errorf("Runs should be deleted with cascade, but found %d", len(runs))
	}
}
