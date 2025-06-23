package utils

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

type ErrorCode string

const (
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrInternal     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Status  int       `json:"-"`
	Err     error     `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("API error [%s]: %s - %s (original: %v)", e.Code, e.Message, e.Details, e.Err)
	}
	return fmt.Sprintf("API error [%s]: %s - %s", e.Code, e.Message, e.Details)
}

func NewAPIError(code ErrorCode, message string, err error) *APIError {
	apiErr := &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}

	switch code {
	case ErrBadRequest:
		apiErr.Status = iris.StatusBadRequest
	case ErrUnauthorized:
		apiErr.Status = iris.StatusUnauthorized
	case ErrForbidden:
		apiErr.Status = iris.StatusForbidden
	case ErrNotFound:
		apiErr.Status = iris.StatusNotFound
	case ErrValidation:
		apiErr.Status = iris.StatusBadRequest
	default:
		apiErr.Status = iris.StatusInternalServerError
		apiErr.Details = "An unexpected error occurred"
	}

	if err != nil {
		apiErr.Details = err.Error()
	}

	return apiErr
}
