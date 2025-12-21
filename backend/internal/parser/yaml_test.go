package parser

import (
	"testing"

	"gantry/internal/models"
)

func TestParse_ValidWorkflow(t *testing.T) {
	yaml := `
name: Test Workflow
on:
  push:
    branches:
      - main
jobs:
  lint:
    runs-on: alpine
    steps:
      - name: Lint code
        run: echo "linting"
  test:
    runs-on: ubuntu
    steps:
      - name: Run tests
        run: echo "testing"
`

	p := NewParser()
	wf, err := p.Parse([]byte(yaml))

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if wf.Name != "Test Workflow" {
		t.Errorf("Expected name 'Test Workflow', got '%s'", wf.Name)
	}

	if len(wf.Jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(wf.Jobs))
	}

	if len(wf.JobOrder) != 2 {
		t.Errorf("Expected job order length 2, got %d", len(wf.JobOrder))
	}

	// Verify job order is preserved
	if wf.JobOrder[0] != "lint" || wf.JobOrder[1] != "test" {
		t.Errorf("Expected job order [lint, test], got %v", wf.JobOrder)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	yaml := `
name: Invalid
invalid yaml syntax {
`

	p := NewParser()
	_, err := p.Parse([]byte(yaml))

	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestValidate_ValidWorkflow(t *testing.T) {
	wf := &models.Workflow{
		Name: "Valid Workflow",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps: []models.Step{
					{Name: "Step 1", Run: "echo test"},
				},
			},
		},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err != nil {
		t.Errorf("Expected no error for valid workflow, got: %v", err)
	}
}

func TestValidate_MissingName(t *testing.T) {
	wf := &models.Workflow{
		Name: "",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps:  []models.Step{{Name: "Step 1", Run: "echo test"}},
			},
		},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err == nil {
		t.Error("Expected error for missing name, got nil")
	}
}

func TestValidate_NoJobs(t *testing.T) {
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err == nil {
		t.Error("Expected error for no jobs, got nil")
	}
}

func TestValidate_NoSteps(t *testing.T) {
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps:  []models.Step{},
			},
		},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err == nil {
		t.Error("Expected error for job with no steps, got nil")
	}
}

func TestValidate_StepMissingName(t *testing.T) {
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps: []models.Step{
					{Name: "", Run: "echo test"},
				},
			},
		},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err == nil {
		t.Error("Expected error for step with no name, got nil")
	}
}

func TestValidate_StepMissingRun(t *testing.T) {
	wf := &models.Workflow{
		Name: "Test",
		Jobs: map[string]models.Job{
			"test": {
				RunsOn: "ubuntu",
				Steps: []models.Step{
					{Name: "Test", Run: ""},
				},
			},
		},
	}

	p := NewParser()
	err := p.Validate(wf)

	if err == nil {
		t.Error("Expected error for step with no run command, got nil")
	}
}
