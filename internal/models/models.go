package models

import "net/http"

type Status string

const (
	ErrInvalidInput Status = "INVALID_INPUT"
	ErrNotFound     Status = "NOT_FOUND"
	ErrUnauthorized Status = "UNAUTHORIZED"
	ErrInternal     Status = "INTERNAL_ERROR"

	StatusOK        Status = "STATUS_OK"
	StatusCreated   Status = "STATUS_CREATED"
	StatusModified  Status = "STATUS_MODIFIED"
	StatusAccepted  Status = "STATUS_ACCEPTED"
	StatusNoContent Status = "STATUS_NOCONTENT"
)

type APIResponse struct {
	Status     Status      `json:"status"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
}

func (e *APIResponse) ToHTTPStatus() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}

	switch e.Status {
	case StatusCreated:
		e.StatusCode = http.StatusCreated
	case StatusOK:
		e.StatusCode = http.StatusOK
	case StatusAccepted:
		e.StatusCode = http.StatusAccepted
	case ErrInvalidInput:
		e.StatusCode = http.StatusBadRequest
	case ErrNotFound:
		e.StatusCode = http.StatusNotFound
	case ErrUnauthorized:
		e.StatusCode = http.StatusUnauthorized
	case ErrInternal:
		e.StatusCode = http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
	return e.StatusCode
}
