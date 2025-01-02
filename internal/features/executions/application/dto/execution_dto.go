package dto

import (
	"faas/internal/features/executions/domain/entity"
	"time"
)

type CreateExecutionRequest struct {
	FunctionID string `json:"function_id" binding:"required"`
	Input      string `json:"input"`
	//Input struct {
	//	DirectInputs map[string]interface{} `json:"direct_inputs,omitempty"`
	//	ObjectInputs map[string]string      `json:"object_inputs,omitempty"`
	//	Secrets      []string               `json:"secrets,omitempty"` // Lista de nombres de secrets a usar
	//} `json:"input" validate:"required"`
}

type ExecutionResponse struct {
	ID          string     `json:"id"`
	FunctionID  string     `json:"function_id"`
	Status      string     `json:"status"`
	Input       string     `json:"input"`
	Output      string     `json:"output,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func NewExecutionResponse(execution *entity.Execution) *ExecutionResponse {
	return &ExecutionResponse{
		ID:          execution.ID,
		FunctionID:  execution.FunctionID,
		Status:      string(execution.Status),
		Input:       execution.Input,
		Output:      execution.Output,
		Error:       execution.Error,
		CreatedAt:   execution.CreatedAt,
		StartedAt:   execution.StartedAt,
		CompletedAt: execution.CompletedAt,
	}
}
