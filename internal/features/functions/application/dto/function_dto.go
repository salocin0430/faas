package dto

import "faas/internal/features/functions/domain/entity"

type CreateFunctionRequest struct {
	Name        string `json:"name" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required"`
	Description string `json:"description"`
}

type FunctionResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	UserID   string `json:"user_id"`
}

func NewFunctionResponse(function *entity.Function) *FunctionResponse {
	return &FunctionResponse{
		ID:       function.ID,
		Name:     function.Name,
		ImageURL: function.ImageURL,
		UserID:   function.UserID,
	}
}
