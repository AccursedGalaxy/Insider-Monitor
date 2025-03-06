package solana

import (
	"testing"
)

func TestIsValidAddress(t *testing.T) {
	// Test valid addresses
	validAddresses := []string{
		"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
		"DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF",
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
	}

	for _, addr := range validAddresses {
		if !IsValidAddress(addr) {
			t.Errorf("Expected %s to be a valid address", addr)
		}
	}

	// Test invalid addresses - note that Solana addresses can vary in length
	// so we focus on truly invalid characters/formats
	invalidAddresses := []string{
		"",
		"invalid",
		"123",
		"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVF!",  // Invalid character
		"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr=", // Invalid character
		"Not-A-Valid-Address",
	}

	for _, addr := range invalidAddresses {
		if IsValidAddress(addr) {
			t.Errorf("Expected %s to be an invalid address", addr)
		}
	}
}

func TestFormatTokenAmount(t *testing.T) {
	testCases := []struct {
		amount   uint64
		decimals uint8
		expected string
	}{
		{100, 0, "100"},
		{100, 1, "10.0"},
		{100, 2, "1.00"},
		{123456789, 6, "123.456789"},
		{1000000, 6, "1.000000"},
		{1, 6, "0.000001"},
		{0, 6, "0.000000"},
		{101, 2, "1.01"},
		{199, 2, "1.99"},
	}

	for _, tc := range testCases {
		result := FormatTokenAmount(tc.amount, tc.decimals)
		if result != tc.expected {
			t.Errorf("FormatTokenAmount(%d, %d) = %s, expected %s",
				tc.amount, tc.decimals, result, tc.expected)
		}
	}
}

func TestTokenToDecimal(t *testing.T) {
	testCases := []struct {
		amount   uint64
		decimals uint8
		expected float64
	}{
		{100, 0, 100.0},
		{100, 1, 10.0},
		{100, 2, 1.0},
		{123456789, 6, 123.456789},
		{1000000, 6, 1.0},
		{1, 6, 0.000001},
		{0, 6, 0.0},
		{101, 2, 1.01},
		{199, 2, 1.99},
	}

	for _, tc := range testCases {
		result := TokenToDecimal(tc.amount, tc.decimals)
		if result != tc.expected {
			t.Errorf("TokenToDecimal(%d, %d) = %f, expected %f",
				tc.amount, tc.decimals, result, tc.expected)
		}
	}
}

func TestPublicKeyFromBase58(t *testing.T) {
	validAddress := "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr"

	// Test valid address
	pubkey, err := PublicKeyFromBase58(validAddress)
	if err != nil {
		t.Errorf("Expected no error for valid address, got %v", err)
	}

	// The string representation should match the original
	if pubkey.String() != validAddress {
		t.Errorf("Expected pubkey.String() to be %s, got %s", validAddress, pubkey.String())
	}

	// Test invalid address
	_, err = PublicKeyFromBase58("invalid-address")
	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
}
