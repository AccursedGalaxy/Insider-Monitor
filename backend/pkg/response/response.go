package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response is a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Err     *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents an API error
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains metadata about the response
type Meta struct {
	Timestamp   time.Time              `json:"timestamp"`
	RequestID   string                 `json:"request_id,omitempty"`
	Count       int                    `json:"count,omitempty"`
	TotalCount  int                    `json:"total_count,omitempty"`
	Page        int                    `json:"page,omitempty"`
	TotalPages  int                    `json:"total_pages,omitempty"`
	NextPage    string                 `json:"next_page,omitempty"`
	PrevPage    string                 `json:"prev_page,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	ProcessTime string                 `json:"process_time,omitempty"`
}

// Common error codes
const (
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrMethodNotAllowed    = "METHOD_NOT_ALLOWED"
	ErrConflict            = "CONFLICT"
	ErrTooManyRequests     = "TOO_MANY_REQUESTS"
	ErrInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
	ErrValidation          = "VALIDATION_ERROR"
)

// Helper functions for common responses

// Success returns a successful response
func Success(c *gin.Context, data interface{}) {
	resp := Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now(),
			RequestID: c.GetString("RequestID"),
		},
	}

	c.JSON(http.StatusOK, resp)
}

// ErrorResponse returns an error response
func ErrorResponse(c *gin.Context, statusCode int, errorCode, message, details string) {
	resp := Response{
		Success: false,
		Err: &ErrorInfo{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
		Meta: &Meta{
			Timestamp: time.Now(),
			RequestID: c.GetString("RequestID"),
		},
	}

	c.JSON(statusCode, resp)
}

// SuccessWithMeta returns a successful response with metadata
func SuccessWithMeta(c *gin.Context, data interface{}, meta Meta) {
	meta.Timestamp = time.Now()
	if meta.RequestID == "" {
		meta.RequestID = c.GetString("RequestID")
	}

	resp := Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	}

	c.JSON(http.StatusOK, resp)
}

// SuccessWithPagination returns a successful response with pagination metadata
func SuccessWithPagination(c *gin.Context, data interface{}, count, totalCount, page, pageSize int) {
	totalPages := totalCount / pageSize
	if totalCount%pageSize > 0 {
		totalPages++
	}

	meta := Meta{
		Timestamp:  time.Now(),
		RequestID:  c.GetString("RequestID"),
		Count:      count,
		TotalCount: totalCount,
		Page:       page,
		TotalPages: totalPages,
	}

	// Add next/prev page URLs if appropriate
	baseURL := c.Request.URL.Path + "?"
	query := c.Request.URL.Query()

	if page < totalPages {
		query.Set("page", string(rune(page+1)))
		meta.NextPage = baseURL + query.Encode()
	}

	if page > 1 {
		query.Set("page", string(rune(page-1)))
		meta.PrevPage = baseURL + query.Encode()
	}

	resp := Response{
		Success: true,
		Data:    data,
		Meta:    &meta,
	}

	c.JSON(http.StatusOK, resp)
}

// BadRequest returns a 400 Bad Request error
func BadRequest(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusBadRequest, ErrBadRequest, message, details)
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusUnauthorized, ErrUnauthorized, message, details)
}

// Forbidden returns a 403 Forbidden error
func Forbidden(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusForbidden, ErrForbidden, message, details)
}

// NotFound returns a 404 Not Found error
func NotFound(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusNotFound, ErrNotFound, message, details)
}

// MethodNotAllowed returns a 405 Method Not Allowed error
func MethodNotAllowed(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusMethodNotAllowed, ErrMethodNotAllowed, message, details)
}

// Conflict returns a 409 Conflict error
func Conflict(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusConflict, ErrConflict, message, details)
}

// TooManyRequests returns a 429 Too Many Requests error
func TooManyRequests(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusTooManyRequests, ErrTooManyRequests, message, details)
}

// InternalServerError returns a 500 Internal Server Error
func InternalServerError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusInternalServerError, ErrInternalServerError, message, details)
}

// ServiceUnavailable returns a 503 Service Unavailable error
func ServiceUnavailable(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusServiceUnavailable, ErrServiceUnavailable, message, details)
}
