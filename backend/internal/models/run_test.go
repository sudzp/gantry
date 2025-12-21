package models

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWorkflowRun_UpdateJob(t *testing.T) {
	run := &WorkflowRun{
		ID:       "run-1",
		Jobs:     make(map[string]Job),
		JobOrder: []string{"test"},
	}

	job := Job{
		RunsOn: "ubuntu",
		Status: "running",
		Steps:  []Step{{Name: "Test", Run: "echo test"}},
	}

	run.UpdateJob("test", job)

	retrieved, exists := run.GetJob("test")
	if !exists {
		t.Error("Expected job to exist after update")
	}

	if retrieved.Status != "running" {
		t.Errorf("Expected status 'running', got '%s'", retrieved.Status)
	}
}

func TestWorkflowRun_GetJob(t *testing.T) {
	run := &WorkflowRun{
		ID:   "run-1",
		Jobs: make(map[string]Job),
	}

	job := Job{RunsOn: "ubuntu", Status: "success"}
	run.Jobs["test"] = job

	retrieved, exists := run.GetJob("test")
	if !exists {
		t.Error("Expected job to exist")
	}

	if retrieved.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", retrieved.Status)
	}

	// Test non-existent job
	_, exists = run.GetJob("nonexistent")
	if exists {
		t.Error("Expected non-existent job to return false")
	}
}

func TestWorkflowRun_SetStatus(t *testing.T) {
	run := &WorkflowRun{
		ID:     "run-1",
		Status: "running",
		Jobs:   make(map[string]Job),
	}

	run.SetStatus("success")

	if run.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", run.Status)
	}
}

func TestWorkflowRun_Complete(t *testing.T) {
	run := &WorkflowRun{
		ID:   "run-1",
		Jobs: make(map[string]Job),
	}

	if run.CompletedAt != nil {
		t.Error("Expected CompletedAt to be nil initially")
	}

	run.Complete()

	if run.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set after Complete()")
	}

	// Check that completed time is recent
	if time.Since(*run.CompletedAt) > time.Second {
		t.Error("CompletedAt should be recent")
	}
}

func TestWorkflowRun_Clone(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(5 * time.Minute)

	original := &WorkflowRun{
		ID:           "run-1",
		WorkflowName: "Test",
		Status:       "success",
		Jobs: map[string]Job{
			"test": {RunsOn: "ubuntu", Status: "success"},
		},
		JobOrder:    []string{"test"},
		StartedAt:   startTime,
		CompletedAt: &endTime,
	}

	clone := original.Clone()

	// Verify clone has same data
	if clone.ID != original.ID {
		t.Errorf("Expected ID '%s', got '%s'", original.ID, clone.ID)
	}

	if clone.Status != original.Status {
		t.Errorf("Expected status '%s', got '%s'", original.Status, clone.Status)
	}

	if len(clone.Jobs) != len(original.Jobs) {
		t.Errorf("Expected %d jobs, got %d", len(original.Jobs), len(clone.Jobs))
	}

	// Verify clone is independent
	clone.Status = "failed"
	if original.Status == "failed" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestWorkflowRun_ThreadSafety(t *testing.T) {
	run := &WorkflowRun{
		ID:       "run-1",
		Jobs:     make(map[string]Job),
		JobOrder: []string{"job1", "job2", "job3"},
	}

	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			job := Job{
				RunsOn: "ubuntu",
				Status: "running",
			}
			run.UpdateJob(fmt.Sprintf("job-%d", id), job)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = run.Clone()
		}()
	}

	wg.Wait()

	// Verify no race conditions occurred
	if len(run.Jobs) < 10 {
		t.Errorf("Expected at least 10 jobs, got %d", len(run.Jobs))
	}
}

func TestWorkflowRun_SetStatusConcurrency(t *testing.T) {
	run := &WorkflowRun{
		ID:     "run-1",
		Status: "running",
		Jobs:   make(map[string]Job),
	}

	var wg sync.WaitGroup

	// Multiple goroutines trying to set status
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(status string) {
			defer wg.Done()
			run.SetStatus(status)
		}(fmt.Sprintf("status-%d", i))
	}

	wg.Wait()

	// Should have one of the statuses (no race condition crash)
	if run.Status == "" {
		t.Error("Status should be set")
	}
}
