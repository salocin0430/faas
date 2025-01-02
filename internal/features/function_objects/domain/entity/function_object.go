package entity

import "time"

type FunctionObject struct {
	ID          string    `json:"id"`
	FunctionID  string    `json:"function_id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
