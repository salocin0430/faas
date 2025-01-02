package dto

type CreateSecretRequest struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type UpdateSecretRequest struct {
	Value string `json:"value" validate:"required"`
}
