package dto

type CreateObjectRequest struct {
	FunctionID string `form:"function_id"`
	Name       string `form:"name"`
}

type ListObjectsRequest struct {
	FunctionID string `uri:"function_id"`
}
