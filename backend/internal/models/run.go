// Package models defines the data structures for workflows and runs
package models

import (
	"sync"
	"time"
)

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

// UpdateJob safely updates a job in the run
func (r *WorkflowRun) UpdateJob(name string, job Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Jobs[name] = job
}

// GetJob safely retrieves a job from the run
func (r *WorkflowRun) GetJob(name string) (Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, exists := r.Jobs[name]
	return job, exists
}

// SetStatus safely sets the run status
func (r *WorkflowRun) SetStatus(status string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Status = status
}

// Complete marks the run as completed
func (r *WorkflowRun) Complete() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.CompletedAt = &now
}

// Clone creates a safe copy for reading
func (r *WorkflowRun) Clone() *WorkflowRun {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := &WorkflowRun{
		ID:           r.ID,
		WorkflowName: r.WorkflowName,
		Status:       r.Status,
		Jobs:         make(map[string]Job),
		JobOrder:     make([]string, len(r.JobOrder)),
		StartedAt:    r.StartedAt,
		CompletedAt:  r.CompletedAt,
	}

	for k, v := range r.Jobs {
		clone.Jobs[k] = v
	}
	copy(clone.JobOrder, r.JobOrder)

	return clone
}
