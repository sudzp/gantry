package storage

import (
	"fmt"
	"sync"

	"gantry/internal/models"
)

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	workflows    map[string]*models.Workflow
	workflowRuns map[string]*models.WorkflowRun
	mu           sync.RWMutex
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		workflows:    make(map[string]*models.Workflow),
		workflowRuns: make(map[string]*models.WorkflowRun),
	}
}

// SaveWorkflow saves a workflow
func (s *MemoryStorage) SaveWorkflow(wf *models.Workflow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workflows[wf.Name] = wf
	return nil
}

// GetWorkflow retrieves a workflow by name
func (s *MemoryStorage) GetWorkflow(name string) (*models.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wf, exists := s.workflows[name]
	if !exists {
		return nil, fmt.Errorf("workflow '%s' not found", name)
	}
	return wf, nil
}

// ListWorkflows returns all workflows
func (s *MemoryStorage) ListWorkflows() ([]*models.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workflows := make([]*models.Workflow, 0, len(s.workflows))
	for _, wf := range s.workflows {
		workflows = append(workflows, wf)
	}
	return workflows, nil
}

// DeleteWorkflow deletes a workflow
func (s *MemoryStorage) DeleteWorkflow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workflows[name]; !exists {
		return fmt.Errorf("workflow '%s' not found", name)
	}
	delete(s.workflows, name)
	return nil
}

// SaveRun saves a workflow run
func (s *MemoryStorage) SaveRun(run *models.WorkflowRun) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workflowRuns[run.ID] = run
	return nil
}

// GetRun retrieves a run by ID
func (s *MemoryStorage) GetRun(id string) (*models.WorkflowRun, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	run, exists := s.workflowRuns[id]
	if !exists {
		return nil, fmt.Errorf("run '%s' not found", id)
	}
	return run.Clone(), nil
}

// ListRuns returns all runs
func (s *MemoryStorage) ListRuns() ([]*models.WorkflowRun, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	runs := make([]*models.WorkflowRun, 0, len(s.workflowRuns))
	for _, run := range s.workflowRuns {
		runs = append(runs, run.Clone())
	}
	return runs, nil
}

// UpdateRun updates an existing run
func (s *MemoryStorage) UpdateRun(run *models.WorkflowRun) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workflowRuns[run.ID]; !exists {
		return fmt.Errorf("run '%s' not found", run.ID)
	}
	s.workflowRuns[run.ID] = run
	return nil
}

// DeleteRunsByWorkflow deletes all runs for a workflow
func (s *MemoryStorage) DeleteRunsByWorkflow(workflowName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, run := range s.workflowRuns {
		if run.WorkflowName == workflowName {
			delete(s.workflowRuns, id)
		}
	}
	return nil
}
