// Purpose: Standard API response envelope — success/error helpers
// Dev must implement: consistent JSON response format across all endpoints
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is the standard success envelope
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the standard error envelope
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents a field-level validation failure
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// OK sends a 200 JSON response
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

// Created sends a 201 JSON response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

// BadRequest sends a 400 validation error response
func BadRequest(c *gin.Context, details []ValidationError) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "validation_error",
		Details: details,
	})
}

// Unauthorized sends a 401 response
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
}

// Forbidden sends a 403 response
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
}

// NotFound sends a 404 response
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "not_found"})
}

// InternalServerError sends a 500 response
func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_server_error"})
}
