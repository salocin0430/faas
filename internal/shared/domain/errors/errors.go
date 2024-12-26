package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}
