package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type Code string

const (
	CodeNotFound     Code = "NOT_FOUND"
	CodeInvalidInput Code = "INVALID_INPUT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeConflict     Code = "CONFLICT"
	CodeInternal     Code = "INTERNAL"
	CodeUnavailable  Code = "UNAVAILABLE"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NotFound(message string, err error) *AppError {
	return New(CodeNotFound, message, err)
}

func Invalid(message string, err error) *AppError {
	return New(CodeInvalidInput, message, err)
}

func Internal(message string, err error) *AppError {
	return New(CodeInternal, message, err)
}

func HTTPStatus(err error) int {
	var ae *AppError
	if errors.As(err, &ae) {
		switch ae.Code {
		case CodeNotFound:
			return http.StatusNotFound
		case CodeInvalidInput:
			return http.StatusBadRequest
		case CodeUnauthorized:
			return http.StatusUnauthorized
		case CodeForbidden:
			return http.StatusForbidden
		case CodeConflict:
			return http.StatusConflict
		case CodeUnavailable:
			return http.StatusServiceUnavailable
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

func ClientMessage(err error) string {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Message
	}
	return "internal server error"
}
