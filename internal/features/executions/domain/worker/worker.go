package worker

import (
	"context"
)

// Worker representa el dominio del worker
type Worker interface {
	Start(ctx context.Context) error
	ProcessTask(ctx context.Context, task Task) error
	Stop() error
}

// Task representa una tarea a ejecutar
type Task struct {
	ID             string         `json:"id"`
	FunctionConfig FunctionConfig `json:"function_config"`
}

// FunctionConfig representa la configuración de una función
type FunctionConfig struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Environment map[string]string `json:"env"`
}

// TaskResult representa el resultado de una tarea
type TaskResult struct {
	TaskID string `json:"task_id"`
	Result string `json:"result"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
