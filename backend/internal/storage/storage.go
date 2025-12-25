package storage

import "gantry/internal/models"

// Storage defines the interface for workflow and run storage
type Storage interface {
	// Workflow operations
	SaveWorkflow(wf *models.Workflow) error
	GetWorkflow(name string) (*models.Workflow, error)
	ListWorkflows() ([]*models.Workflow, error)
	DeleteWorkflow(name string) error

	// Run operations
	SaveRun(run *models.WorkflowRun) error
	GetRun(id string) (*models.WorkflowRun, error)
	ListRuns() ([]*models.WorkflowRun, error)
	UpdateRun(run *models.WorkflowRun) error
	DeleteRunsByWorkflow(workflowName string) error
}
