package utils

import (
	"fiber-clean-transaction/pkg/validation"
	"fmt"
	"net/http"
)

// Error codes constants
const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
)

// Domain Error
type DomainError struct {
	Code       string                     `json:"code"`
	Message    string                     `json:"message"`
	StatusCode int                        `json:"status_code"`
	Err        error                      `json:"-"`
	Errors     []validation.DetailedError `json:"errors"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Error constructors
func NotFound(message string) *DomainError {
	return &DomainError{
		Code:       ErrCodeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func BadRequest(message string) *DomainError {
	return &DomainError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Internal(message string, err error) *DomainError {
	return &DomainError{
		Code:       ErrCodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func Unauthorized(message string) *DomainError {
	return &DomainError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// func ErrorValidation(errors []ValidationError) *DomainError {
// 	return &DomainError{
// 		Code      : ErrCodeBadRequest,
// 		Message   : "Validation failed",
// 		StatusCode: http.StatusUnprocessableEntity,
// 		Errors    : errors,
// 	}
// }

func ErrorValidation(errors []validation.DetailedError) *DomainError {
	return &DomainError{
		Code:       ErrCodeBadRequest,
		Message:    "Validation failed",
		StatusCode: http.StatusUnprocessableEntity,
		Errors:     errors,
	}
}
