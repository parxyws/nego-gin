package dto

import "github.com/parxyws/nego-gin/internal/shared/domain"

type ApiResponse[T any] struct {
	Success  bool            `json:"success"`
	Message  string          `json:"message"`
	Data     T               `json:"data,omitempty"`
	Error    error           `json:"error,omitempty"`
	Metadata domain.Metadata `json:"metadata,omitempty"`
}
