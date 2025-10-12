package models

import (
	"net/http"
)

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

// AI data models
type TranslateApiBody struct {
	TextToTranslate string `json:"text_to_translate" binding:"required"`
	Language        string `json:"language" binding:"required"`
}

type InventoryApiBody struct {
	InventoryInput string `json:"inventory_input" binding:"required"`
}

type InventoryData struct {
	Item                      string  `json:"item"`
	Quantity                  float32 `json:"quantity"`
	Unit                      string  `json:"unit"`
	WholesalePricePerQuantity float32 `json:"wholesale_price_per_quantity"`
	TotalCostOfProduct        float32 `json:"total_cost_of_product"`
}

// AI data models

// global api response models
type APIResponse struct {
	Status     Status      `json:"status"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
}

// global api response models

// transcription models
type TranscriptionResponse struct {
	Results struct {
		Transcripts []Transcript `json:"transcripts"`
	} `json:"results"`
}

type Transcript struct {
	Transcript string `json:"transcript"`
}

// transcription models

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
	default:
		e.StatusCode = http.StatusInternalServerError
	}
	return e.StatusCode
}
