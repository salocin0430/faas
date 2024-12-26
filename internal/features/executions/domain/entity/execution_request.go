package entity

import "time"

type ExecutionRequest struct {
	ID       string `json:"id"`
	ImageURL string `json:"image_url"`
	Input    string `json:"input"`
	UserID   string `json:"user_id"`
}

type ExecutionResult struct {
	ExecutionID string    `json:"execution_id"`
	Status      string    `json:"status"`
	Output      string    `json:"output,omitempty"`
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}
