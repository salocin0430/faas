package entity

import (
	"time"

	"github.com/google/uuid"
)

type Function struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewFunction(name, imageURL, userID string) *Function {
	now := time.Now()
	return &Function{
		ID:        uuid.New().String(),
		Name:      name,
		ImageURL:  imageURL,
		UserID:    userID,
		CreatedAt: now,
	}
}
