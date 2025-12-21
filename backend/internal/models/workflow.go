package models

import "time"

// Workflow defines the CI/CD pipeline structure
type Workflow struct {
	Name     string         `yaml:"name" json:"name"`
	On       TriggerConfig  `yaml:"on" json:"on"`
	Jobs     map[string]Job `yaml:"jobs" json:"jobs"`
	JobOrder []string       `json:"job_order"` // Preserve YAML order
}

// TriggerConfig defines when the workflow triggers
type TriggerConfig struct {
	Push PushConfig `yaml:"push"`
}

// PushConfig defines push trigger configuration
type PushConfig struct {
	Branches []string `yaml:"branches"`
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

// JobResult contains the result of job execution
type JobResult struct {
	Output string
	Error  error
}
