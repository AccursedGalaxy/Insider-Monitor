package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetWallets(t *testing.T) {
	// Set up the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/wallets", GetWallets)

	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/api/v1/wallets", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(recorder, req)

	// Check the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Parse the response body
	var response []map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check that we got at least one wallet
	if len(response) == 0 {
		t.Error("Expected at least one wallet in response")
	}

	// Check that each wallet has the expected fields
	for _, wallet := range response {
		// Check required fields
		requiredFields := []string{"address", "label", "last_scanned", "token_count"}
		for _, field := range requiredFields {
			if _, ok := wallet[field]; !ok {
				t.Errorf("Wallet missing required field: %s", field)
			}
		}
	}
}

func TestGetWalletByAddress(t *testing.T) {
	// Set up the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/wallets/:address", GetWalletByAddress)

	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/api/v1/wallets/55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(recorder, req)

	// Check the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Parse the response body
	var wallet map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &wallet); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"address", "label", "last_scanned", "token_count"}
	for _, field := range requiredFields {
		if _, ok := wallet[field]; !ok {
			t.Errorf("Wallet missing required field: %s", field)
		}
	}

	// Check address
	address, ok := wallet["address"].(string)
	if !ok {
		t.Error("Address field is not a string")
	} else if address != "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr" {
		t.Errorf("Expected address %s, got %s", "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", address)
	}
}

func TestGetWalletTokens(t *testing.T) {
	// Set up the Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/wallets/:address/tokens", GetWalletTokens)

	// Create a request to the endpoint
	req, err := http.NewRequest("GET", "/api/v1/wallets/55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr/tokens", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(recorder, req)

	// Check the status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check address
	address, ok := response["address"].(string)
	if !ok {
		t.Error("Address field is not a string")
	} else if address != "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr" {
		t.Errorf("Expected address %s, got %s", "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", address)
	}

	// Check tokens
	tokens, ok := response["tokens"].([]interface{})
	if !ok {
		t.Error("Tokens field is not an array")
	} else if len(tokens) == 0 {
		t.Error("Expected at least one token in response")
	} else {
		// Check that each token has the expected fields
		token := tokens[0].(map[string]interface{})
		requiredFields := []string{"mint", "symbol", "balance", "decimals", "last_updated"}
		for _, field := range requiredFields {
			if _, ok := token[field]; !ok {
				t.Errorf("Token missing required field: %s", field)
			}
		}
	}
}
