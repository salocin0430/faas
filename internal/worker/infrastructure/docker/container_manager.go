package docker

import (
	"bytes"
	"context"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/worker/domain/ports"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerContainerManager struct {
	client       *client.Client
	functionRepo ports.FunctionRepository
}

func NewContainerManager(functionRepo ports.FunctionRepository) (ports.ContainerManager, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.46"),
	)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManager{
		client:       cli,
		functionRepo: functionRepo,
	}, nil
}

func (m *DockerContainerManager) RunFunction(ctx context.Context, execution *entity.Execution) (string, error) {
	// Get function from repository
	function, err := m.functionRepo.GetByID(ctx, execution.FunctionID)
	if err != nil {
		return "", err
	}

	// Pull image if needed
	reader, err := m.client.ImagePull(ctx, function.ImageURL, image.PullOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)

	// Create container with input as argument
	var cmd []string
	if execution.Input != "" {
		cmd = []string{execution.Input} // Only add input if not empty
	}

	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image: function.ImageURL,
		Cmd:   cmd, // Can be empty or contain input
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start container
	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Create context with timeout
	execTimeout := 5 * time.Minute // or preferred value
	ctx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	// Wait for container to finish
	statusCh, errCh := m.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return "", err
	case <-statusCh:
		log.Printf("Container %s finished execution", resp.ID)
	case <-ctx.Done():
		return "", fmt.Errorf("execution timed out after %v", execTimeout)
	}

	// Get logs
	out, err := m.client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Read output and clean control bytes
	var stdoutBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, io.Discard, out)
	if err != nil {
		return "", err
	}

	// Cleanup
	m.client.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	return stdoutBuf.String(), nil
}

func (m *DockerContainerManager) Stop() error {
	return m.client.Close()
}
