package dto

import (
	"faas/internal/features/function_objects/domain/entity"
	"time"
)

type ObjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewObjectResponse(obj *entity.FunctionObject) *ObjectResponse {
	return &ObjectResponse{
		ID:          obj.ID,
		Name:        obj.Name,
		Size:        obj.Size,
		ContentType: obj.ContentType,
		CreatedAt:   obj.CreatedAt,
	}
}
