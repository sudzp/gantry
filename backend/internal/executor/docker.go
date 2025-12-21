package executor

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"gantry/internal/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerExecutor executes jobs using Docker containers
type DockerExecutor struct {
	client *client.Client
}

// NewDockerExecutor creates a new Docker-based executor
func NewDockerExecutor() (*DockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerExecutor{
		client: cli,
	}, nil
}

// Execute runs a job in a Docker container
func (e *DockerExecutor) Execute(ctx context.Context, jobName string, job models.Job) (string, error) {
	// Use background context for Docker operations to avoid premature cancellation
	// Create separate timeouts for each operation

	// Select image based on runs-on
	image := "ubuntu:latest"
	if job.RunsOn == "alpine" {
		image = "alpine:latest"
	}

	// Build script with step tracking and timestamps
	script := "#!/bin/sh\nset -e\n"
	for i, step := range job.Steps {
		script += fmt.Sprintf("\n# Step %d: %s\n", i+1, step.Name)
		script += fmt.Sprintf("echo '=== [' $(date '+%%Y-%%m-%%d %%H:%%M:%%S') '] Starting: %s ==='\n", step.Name)
		script += step.Run + "\n"
		script += fmt.Sprintf("echo '=== [' $(date '+%%Y-%%m-%%d %%H:%%M:%%S') '] Completed: %s ==='\n", step.Name)
	}

	// Pull image with separate context and timeout
	log.Printf("Pulling image %s...", image)
	pullCtx, pullCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer pullCancel()

	reader, err := e.client.ImagePull(pullCtx, image, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	// Must read the response to completion
	_, err = io.Copy(io.Discard, reader)
	reader.Close()
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	log.Printf("Image %s pulled successfully", image)

	// Create container with separate context
	createCtx, createCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer createCancel()

	resp, err := e.client.ContainerCreate(createCtx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/sh", "-c", script},
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container with separate context
	startCtx, startCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer startCancel()

	if err := e.client.ContainerStart(startCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for completion with longer timeout (use parent context here)
	statusCh, errCh := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			// Get logs even on failure
			logs := e.getContainerLogs(resp.ID)
			e.cleanupContainer(resp.ID)
			return logs, fmt.Errorf("container exited with status %d", status.StatusCode)
		}
	}

	// Get logs
	logs := e.getContainerLogs(resp.ID)

	// Remove container
	e.cleanupContainer(resp.ID)

	return logs, nil
}

// getContainerLogs retrieves logs from a container
func (e *DockerExecutor) getContainerLogs(containerID string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	out, err := e.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return fmt.Sprintf("Failed to get logs: %v", err)
	}
	defer out.Close()

	logs, _ := io.ReadAll(out)
	return string(logs)
}

// cleanupContainer removes a container
func (e *DockerExecutor) cleanupContainer(containerID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	e.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
}

// Cleanup performs any necessary cleanup
func (e *DockerExecutor) Cleanup() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}
