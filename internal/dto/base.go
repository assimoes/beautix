package dto

import "time"

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(page, pageSize int, total int64) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// BaseResponse contains common response fields
type BaseResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Data       []*T                `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// NewListResponse creates a new list response
func NewListResponse[T any](data []*T, pagination *PaginationResponse) *ListResponse[T] {
	return &ListResponse[T]{
		Data:       data,
		Pagination: pagination,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SearchCriteria represents search parameters
type SearchCriteria struct {
	Query    string         `json:"query,omitempty"`
	Filters  map[string]any `json:"filters,omitempty"`
	SortBy   string                 `json:"sort_by,omitempty"`
	SortDesc bool                   `json:"sort_desc,omitempty"`
}

// DefaultPagination returns default pagination values
func DefaultPagination() *PaginationRequest {
	return &PaginationRequest{
		Page:     1,
		PageSize: 20,
	}
}