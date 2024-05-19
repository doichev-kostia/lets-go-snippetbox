package main

import (
	"fmt"
	"net/http"
)

const (
	ErrorInvalidArgument    = "INVALID_ARGUMENT"    // The client specified an invalid argument regardless of the state of the system.
	ErrorFailedPrecondition = "FAILED_PRECONDITION" // The operation was rejected because the system is not in a state required for the operation's execution.
	ErrorNotFound           = "NOT_FOUND"           // The requested entity was not found.
	ErrorAlreadyExists      = "ALREADY_EXISTS"      // The entity that a client tried to create already exists.
	ErrorUnauthenticated    = "UNAUTHENTICATED"     // The caller does not have valid authentication credentials for the operation.
	ErrorPermissionDenied   = "PERMISSION_DENIED"   // The caller does not have permission to execute the specified operation.
	ErrorTooManyRequests    = "TOO_MANY_REQUESTS"   // The caller has exhausted their rate limit or quota
	ErrorInternal           = "INTERNAL"            // The part of the underlying system is broken
	ErrorUnknown            = "UNKNOWN"             // When the application doesn't know how to handle the caught error
	ErrorUnavailable        = "UNAVAILABLE"         // The service is currently unavailable. Can be retried with a backoff.
)

var StatusCodeMap = map[string]int{
	ErrorInvalidArgument:    http.StatusBadRequest,
	ErrorFailedPrecondition: http.StatusBadRequest,
	ErrorNotFound:           http.StatusNotFound,
	ErrorAlreadyExists:      http.StatusConflict,
	ErrorUnauthenticated:    http.StatusUnauthorized,
	ErrorPermissionDenied:   http.StatusForbidden,
	ErrorTooManyRequests:    http.StatusTooManyRequests,
	ErrorInternal:           http.StatusInternalServerError,
	ErrorUnknown:            http.StatusInternalServerError,
	ErrorUnavailable:        http.StatusServiceUnavailable,
}

const (
	typeBadRequest = "BAD_REQUEST"
	typeErrorInfo  = "ERROR_INFO"
)

type ErrorInfo struct {
	Type     string                 `json:"@type"`
	Reason   string                 `json:"reason"`
	Metadata map[string]interface{} `json:"metadata"`
}

type FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

type BadRequest struct {
	Type            string           `json:"@type"`
	FieldViolations []FieldViolation `json:"fieldViolations"`
}

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// Details provide more context to an error. The predefined structs are ErrorInfo and BadRequest
	Details []interface{} `json:"details"`
}

func (error ApiError) Error() string {
	return fmt.Sprintf("ApiError(%s): %s", error.Code, error.Message)
}

func NewApiError(status string, err error, details []interface{}) ApiError {
	if details == nil {
		details = make([]interface{}, 0)
	}
	return ApiError{
		status,
		err.Error(),
		details,
	}
}

func NewBadRequestError(message string, violations []FieldViolation) ApiError {
	details := make([]interface{}, 0)
	if violations != nil {
		details = append(details, BadRequest{typeBadRequest, violations})
	}

	return ApiError{
		Code:    ErrorInvalidArgument,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message string, details []interface{}) ApiError {
	if message == "" {
		message = "Not Found"
	}

	if details == nil {
		details = make([]interface{}, 0)
	}

	return ApiError{
		Code:    ErrorNotFound,
		Message: message,
		Details: details,
	}
}

func NewInternalError() ApiError {
	return ApiError{
		Code:    ErrorInternal,
		Message: "Internal Error",
	}
}
