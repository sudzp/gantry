package storage

import (
	"fmt"
	"testing"
	"time"

	"gantry/internal/models"
)

func TestMemoryStorage_SaveAndGetWorkflow(t *testing.T) {
	store := NewMemoryStorage()

	wf := &models.Workflow{
		Name: "Test Workflow",
		Jobs: map[string]models.Job{
			"test": {RunsOn: "ubuntu", Steps: []models.Step{{Name: "Step 1", Run: "echo test"}}},
		},
	}

	// Save
	err := store.SaveWorkflow(wf)
	if err != nil {
		t.Fatalf("Failed to save workflow: %v", err)
	}

	// Get
	retrieved, err := store.GetWorkflow("Test Workflow")
	if err != nil {
		t.Fatalf("Failed to get workflow: %v", err)
	}

	if retrieved.Name != wf.Name {
		t.Errorf("Expected name '%s', got '%s'", wf.Name, retrieved.Name)
	}
}

func TestMemoryStorage_GetNonExistentWorkflow(t *testing.T) {
	store := NewMemoryStorage()

	_, err := store.GetWorkflow("NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent workflow, got nil")
	}
}

func TestMemoryStorage_ListWorkflows(t *testing.T) {
	store := NewMemoryStorage()

	wf1 := &models.Workflow{Name: "Workflow 1", Jobs: map[string]models.Job{}}
	wf2 := &models.Workflow{Name: "Workflow 2", Jobs: map[string]models.Job{}}

	_ = store.SaveWorkflow(wf1)
	_ = store.SaveWorkflow(wf2)

	workflows, err := store.ListWorkflows()
	if err != nil {
		t.Fatalf("Failed to list workflows: %v", err)
	}

	if len(workflows) != 2 {
		t.Errorf("Expected 2 workflows, got %d", len(workflows))
	}
}

func TestMemoryStorage_DeleteWorkflow(t *testing.T) {
	store := NewMemoryStorage()

	wf := &models.Workflow{Name: "Test", Jobs: map[string]models.Job{}}
	_ = store.SaveWorkflow(wf)

	// Delete
	err := store.DeleteWorkflow("Test")
	if err != nil {
		t.Fatalf("Failed to delete workflow: %v", err)
	}

	// Verify deleted
	_, err = store.GetWorkflow("Test")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestMemoryStorage_DeleteNonExistent(t *testing.T) {
	store := NewMemoryStorage()

	err := store.DeleteWorkflow("NonExistent")
	if err == nil {
		t.Error("Expected error when deleting non-existent workflow, got nil")
	}
}

func TestMemoryStorage_SaveAndGetRun(t *testing.T) {
	store := NewMemoryStorage()

	run := &models.WorkflowRun{
		ID:           "run-123",
		WorkflowName: "Test",
		Status:       "running",
		Jobs:         map[string]models.Job{},
		StartedAt:    time.Now(),
	}

	// Save
	err := store.SaveRun(run)
	if err != nil {
		t.Fatalf("Failed to save run: %v", err)
	}

	// Get
	retrieved, err := store.GetRun("run-123")
	if err != nil {
		t.Fatalf("Failed to get run: %v", err)
	}

	if retrieved.ID != run.ID {
		t.Errorf("Expected ID '%s', got '%s'", run.ID, retrieved.ID)
	}
}

func TestMemoryStorage_UpdateRun(t *testing.T) {
	store := NewMemoryStorage()

	run := &models.WorkflowRun{
		ID:           "run-123",
		WorkflowName: "Test",
		Status:       "running",
		Jobs:         map[string]models.Job{},
		StartedAt:    time.Now(),
	}

	_ = store.SaveRun(run)

	// Update status
	run.Status = "success"
	err := store.UpdateRun(run)
	if err != nil {
		t.Fatalf("Failed to update run: %v", err)
	}

	// Verify update
	retrieved, _ := store.GetRun("run-123")
	if retrieved.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", retrieved.Status)
	}
}

func TestMemoryStorage_ListRuns(t *testing.T) {
	store := NewMemoryStorage()

	run1 := &models.WorkflowRun{ID: "run-1", WorkflowName: "Test", Status: "success", Jobs: map[string]models.Job{}, StartedAt: time.Now()}
	run2 := &models.WorkflowRun{ID: "run-2", WorkflowName: "Test", Status: "failed", Jobs: map[string]models.Job{}, StartedAt: time.Now()}

	_ = store.SaveRun(run1)
	_ = store.SaveRun(run2)

	runs, err := store.ListRuns()
	if err != nil {
		t.Fatalf("Failed to list runs: %v", err)
	}

	if len(runs) != 2 {
		t.Errorf("Expected 2 runs, got %d", len(runs))
	}
}

func TestMemoryStorage_Concurrency(t *testing.T) {
	store := NewMemoryStorage()

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			wf := &models.Workflow{
				Name: fmt.Sprintf("Workflow %d", id),
				Jobs: map[string]models.Job{},
			}
			_ = store.SaveWorkflow(wf)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	workflows, _ := store.ListWorkflows()
	if len(workflows) != 10 {
		t.Errorf("Expected 10 workflows, got %d", len(workflows))
	}
}

func TestMemoryStorage_DeleteRunsByWorkflow(t *testing.T) {
	store := NewMemoryStorage()

	// Save workflow
	wf := &models.Workflow{Name: "TestWorkflow", Jobs: map[string]models.Job{}}
	_ = store.SaveWorkflow(wf)

	// Save runs for this workflow
	run1 := &models.WorkflowRun{ID: "run-1", WorkflowName: "TestWorkflow", Status: "success", Jobs: make(map[string]models.Job)}
	run2 := &models.WorkflowRun{ID: "run-2", WorkflowName: "TestWorkflow", Status: "success", Jobs: make(map[string]models.Job)}
	run3 := &models.WorkflowRun{ID: "run-3", WorkflowName: "OtherWorkflow", Status: "success", Jobs: make(map[string]models.Job)}

	_ = store.SaveRun(run1)
	_ = store.SaveRun(run2)
	_ = store.SaveRun(run3)

	// Verify we have 3 runs
	runs, _ := store.ListRuns()
	if len(runs) != 3 {
		t.Fatalf("Expected 3 runs initially, got %d", len(runs))
	}

	// Delete runs for TestWorkflow
	err := store.DeleteRunsByWorkflow("TestWorkflow")
	if err != nil {
		t.Fatalf("Failed to delete runs: %v", err)
	}

	// Verify only run3 remains
	runs, _ = store.ListRuns()
	if len(runs) != 1 {
		t.Errorf("Expected 1 run after deletion, got %d", len(runs))
	}

	if runs[0].ID != "run-3" {
		t.Errorf("Expected run-3 to remain, got %s", runs[0].ID)
	}
}
