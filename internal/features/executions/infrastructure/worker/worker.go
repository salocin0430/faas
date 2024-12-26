package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"githubcom/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	"faas/internal/features/executions/domain/entity"
)

type Worker struct {
	id           int
	natsJS       nats.JetStreamContext
	dockerClient *client.Client
}

func NewWorker(id int, js nats.JetStreamContext, docker *client.Client) *Worker {
	return &Worker{
		id:           id,
		natsJS:       js,
		dockerClient: docker,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	sub, err := w.natsJS.PullSubscribe(
		"executions.pending",
		"workers",
		nats.MaxInflight(1),
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msgs, _ := sub.Fetch(1)
			for _, msg := range msgs {
				go w.processExecution(ctx, msg)
			}
		}
	}
}

func (w *Worker) processExecution(ctx context.Context, msg *nats.Msg) {
	var execution entity.ExecutionRequest
	if err := json.Unmarshal(msg.Data, &execution); err != nil {
		w.reportError(execution.ID, fmt.Errorf("failed to unmarshal execution: %w", err))
		return
	}

	// Crear contenedor
	containerConfig := &container.Config{
		Image: execution.ImageURL,
		Cmd:   []string{execution.Input},
	}

	cont, err := w.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		w.reportError(execution.ID, fmt.Errorf("failed to create container: %w", err))
		return
	}

	// Iniciar contenedor
	if err := w.dockerClient.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		w.reportError(execution.ID, fmt.Errorf("failed to start container: %w", err))
		return
	}

	// Esperar a que termine
	statusCh, errCh := w.dockerClient.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		w.reportError(execution.ID, fmt.Errorf("container wait failed: %w", err))
	case <-statusCh:
		// Obtener logs
		logs, err := w.dockerClient.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			w.reportError(execution.ID, fmt.Errorf("failed to get logs: %w", err))
			return
		}
		defer logs.Close()

		output, err := io.ReadAll(logs)
		if err != nil {
			w.reportError(execution.ID, fmt.Errorf("failed to read logs: %w", err))
			return
		}
		w.reportSuccess(execution.ID, string(output))
	}

	// Limpiar contenedor
	err = w.dockerClient.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Printf("failed to remove container: %v\n", err)
	}
}

func (w *Worker) reportError(executionID string, err error) {
	result := map[string]interface{}{
		"status": "error",
		"error":  err.Error(),
	}
	data, _ := json.Marshal(result)
	w.natsJS.Publish(fmt.Sprintf("executions.%s.result", executionID), data)
}

func (w *Worker) reportSuccess(executionID string, output string) {
	result := map[string]interface{}{
		"status": "success",
		"output": output,
	}
	data, _ := json.Marshal(result)
	w.natsJS.Publish(fmt.Sprintf("executions.%s.result", executionID), data)
}
