package models

// CreateProductResponse represents the response for creating a product
type CreateProductResponse struct {
	Message string `json:"message" example:"Product created successfully"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error" example:"Bad Request"`
}

// DeleteProductResponse represents the response for deleting a product
type DeleteProductResponse struct {
	Message string `json:"message" example:"Product deleted successfully"`
}
