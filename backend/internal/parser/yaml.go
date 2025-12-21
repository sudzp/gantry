package parser

import (
	"bytes"
	"fmt"

	"gantry/internal/models"

	"gopkg.in/yaml.v3"
)

// Parser handles workflow parsing
type Parser struct{}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a YAML workflow file and preserves job order
func (p *Parser) Parse(data []byte) (*models.Workflow, error) {
	// First, parse the raw YAML to preserve key order
	var rawMap map[string]interface{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Now parse into struct
	var wf models.Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	// Extract job order from raw YAML using yaml.v3 Node
	var orderedYAML struct {
		Jobs yaml.Node `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(data, &orderedYAML); err == nil {
		// Extract keys in order from yaml.Node
		if orderedYAML.Jobs.Kind == yaml.MappingNode {
			wf.JobOrder = make([]string, 0, len(orderedYAML.Jobs.Content)/2)
			for i := 0; i < len(orderedYAML.Jobs.Content); i += 2 {
				if orderedYAML.Jobs.Content[i].Value != "" {
					wf.JobOrder = append(wf.JobOrder, orderedYAML.Jobs.Content[i].Value)
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

// Validate validates a workflow
func (p *Parser) Validate(wf *models.Workflow) error {
	if wf.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(wf.Jobs) == 0 {
		return fmt.Errorf("workflow must have at least one job")
	}

	for jobName, job := range wf.Jobs {
		if len(job.Steps) == 0 {
			return fmt.Errorf("job '%s' must have at least one step", jobName)
		}

		for i, step := range job.Steps {
			if step.Name == "" {
				return fmt.Errorf("job '%s' step %d is missing a name", jobName, i+1)
			}
			if step.Run == "" {
				return fmt.Errorf("job '%s' step '%s' is missing run commands", jobName, step.Name)
			}
		}
	}

	return nil
}
