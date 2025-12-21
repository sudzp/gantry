package executor

import (
	"context"
	"gantry/internal/models"
)

// Executor defines the interface for job execution
type Executor interface {
	// Execute runs a single job and returns the output
	Execute(ctx context.Context, jobName string, job models.Job) (string, error)

	// Cleanup performs any necessary cleanup
	Cleanup() error
}

// Config holds executor configuration
type Config struct {
	DockerHost string
	Timeout    int // seconds
}
