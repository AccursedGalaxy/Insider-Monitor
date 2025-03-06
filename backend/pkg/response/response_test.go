package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestContext sets up a test gin context and recorder
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestSuccessResponse(t *testing.T) {
	c, w := setupTestContext()

	// Test data
	data := map[string]string{
		"message": "test message",
	}

	// Call function under test
	Success(c, data)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	assert.True(t, response.Success)
	assert.Nil(t, response.Err)

	// Check data
	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test message", responseData["message"])
}

func TestErrorResponse(t *testing.T) {
	c, w := setupTestContext()

	// Call function under test
	ErrorResponse(c, http.StatusBadRequest, ErrBadRequest, "Invalid input", "Field 'name' is required")

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Parse response
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)

	// Check error fields
	assert.NotNil(t, response.Err)
	assert.Equal(t, ErrBadRequest, response.Err.Code)
	assert.Equal(t, "Invalid input", response.Err.Message)
	assert.Equal(t, "Field 'name' is required", response.Err.Details)
}

func TestSuccessWithMeta(t *testing.T) {
	c, w := setupTestContext()

	// Test data
	data := []string{"item1", "item2"}

	// Meta data
	meta := Meta{
		Timestamp:  time.Now().UTC().Truncate(time.Millisecond),
		RequestID:  "req-123",
		Count:      2,
		TotalCount: 10,
		Page:       1,
		TotalPages: 5,
	}

	// Call function under test
	SuccessWithMeta(c, data, meta)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	assert.True(t, response.Success)
	assert.Nil(t, response.Err)

	// Check meta
	assert.NotNil(t, response.Meta)
	assert.Equal(t, meta.RequestID, response.Meta.RequestID)
	assert.Equal(t, meta.Count, response.Meta.Count)
	assert.Equal(t, meta.TotalCount, response.Meta.TotalCount)
	assert.Equal(t, meta.Page, response.Meta.Page)
	assert.Equal(t, meta.TotalPages, response.Meta.TotalPages)
}

func TestSuccessWithPagination(t *testing.T) {
	c, w := setupTestContext()

	// Test with pagination parameters
	c.Request.URL.RawQuery = "page=2&per_page=5"

	// Test data
	data := []string{"item1", "item2", "item3", "item4", "item5"}

	// Call function under test
	SuccessWithPagination(c, data, 5, 15, 2, 5)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response fields
	assert.True(t, response.Success)
	assert.Nil(t, response.Err)

	// Check meta
	assert.NotNil(t, response.Meta)
	assert.Equal(t, 5, response.Meta.Count)
	assert.Equal(t, 15, response.Meta.TotalCount)
	assert.Equal(t, 2, response.Meta.Page)
	assert.Equal(t, 3, response.Meta.TotalPages)
	assert.NotEmpty(t, response.Meta.NextPage) // Page 3
	assert.NotEmpty(t, response.Meta.PrevPage) // Page 1
}

func TestCommonErrorResponses(t *testing.T) {
	// Test cases for each type of error response
	testCases := []struct {
		name              string
		handler           func(*gin.Context, string, string)
		expectedStatus    int
		expectedErrorCode string
	}{
		{"BadRequest", BadRequest, http.StatusBadRequest, ErrBadRequest},
		{"Unauthorized", Unauthorized, http.StatusUnauthorized, ErrUnauthorized},
		{"Forbidden", Forbidden, http.StatusForbidden, ErrForbidden},
		{"NotFound", NotFound, http.StatusNotFound, ErrNotFound},
		{"MethodNotAllowed", MethodNotAllowed, http.StatusMethodNotAllowed, ErrMethodNotAllowed},
		{"Conflict", Conflict, http.StatusConflict, ErrConflict},
		{"TooManyRequests", TooManyRequests, http.StatusTooManyRequests, ErrTooManyRequests},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError, ErrInternalServerError},
		{"ServiceUnavailable", ServiceUnavailable, http.StatusServiceUnavailable, ErrServiceUnavailable},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestContext()

			// Call function under test
			tc.handler(c, "Error message", "Error details")

			// Assert response status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Parse response
			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check response fields
			assert.False(t, response.Success)
			assert.Nil(t, response.Data)

			// Check error fields
			assert.NotNil(t, response.Err)
			assert.Equal(t, tc.expectedErrorCode, response.Err.Code)
			assert.Equal(t, "Error message", response.Err.Message)
			assert.Equal(t, "Error details", response.Err.Details)
		})
	}
}
