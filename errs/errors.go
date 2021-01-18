package errs

import "net/http"

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

func (appError AppError) AsMessage() *AppError {
	return &AppError{
		Message: appError.Message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Message: message,
		Code: http.StatusConflict,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Message: message,
		Code: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

type ValidationError struct {
	Code    int    `json:",omitempty"`
	InvalidFields []string `json:"invalid_fields"`
}

func NewValidationError(invalidFields []string) *ValidationError {
	return &ValidationError{
		InvalidFields: invalidFields,
		Code:    http.StatusUnprocessableEntity,
	}
}

func (validationError ValidationError) AsMessage() *ValidationError {
	return &ValidationError{
		InvalidFields: validationError.InvalidFields,
	}
}