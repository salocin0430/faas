package entity

import "time"

type ExecutionStatus string

const (
	StatusPending   ExecutionStatus = "pending"
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
)

type Execution struct {
	ID          string          `json:"id"`
	FunctionID  string          `json:"function_id"`
	UserID      string          `json:"user_id"`
	Status      ExecutionStatus `json:"status"`
	Input       string          `json:"input"`
	Output      string          `json:"output,omitempty"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}
