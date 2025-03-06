package monitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWalletMonitor(t *testing.T) {
	tests := []struct {
		name        string
		networkURL  string
		wallets     []string
		shouldError bool
	}{
		{
			name:       "Valid initialization",
			networkURL: "https://api.mainnet-beta.solana.com",
			wallets: []string{
				"DYw8jCTfwHNRJhhmFcbXvVDTqWMEVFBX6ZKUmG5CNSKK", // Example valid wallet
			},
			shouldError: false,
		},
		{
			name:       "Invalid wallet address",
			networkURL: "https://api.mainnet-beta.solana.com",
			wallets: []string{
				"invalid-address",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor, err := NewWalletMonitor(tt.networkURL, tt.wallets, nil)
			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, monitor)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, monitor)
				assert.Equal(t, tt.networkURL, monitor.networkURL)
				assert.Len(t, monitor.wallets, len(tt.wallets))
			}
		})
	}
}

func TestCalculatePercentageChange(t *testing.T) {
	tests := []struct {
		name     string
		old      uint64
		new      uint64
		expected float64
	}{
		{
			name:     "100% increase",
			old:      100,
			new:      200,
			expected: 100.0,
		},
		{
			name:     "50% decrease",
			old:      200,
			new:      100,
			expected: -50.0,
		},
		{
			name:     "New addition",
			old:      0,
			new:      100,
			expected: 100.0,
		},
		{
			name:     "No change",
			old:      100,
			new:      100,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentageChange(tt.old, tt.new)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectChanges(t *testing.T) {
	oldData := map[string]*WalletData{
		"wallet1": {
			WalletAddress: "wallet1",
			TokenAccounts: map[string]TokenAccountInfo{
				"token1": {
					Balance:     1000,
					LastUpdated: time.Now(),
					Symbol:      "TKN1",
					Decimals:    9,
				},
			},
		},
	}

	newData := map[string]*WalletData{
		"wallet1": {
			WalletAddress: "wallet1",
			TokenAccounts: map[string]TokenAccountInfo{
				"token1": {
					Balance:     2000,
					LastUpdated: time.Now(),
					Symbol:      "TKN1",
					Decimals:    9,
				},
			},
		},
	}

	changes := DetectChanges(oldData, newData, 50.0)
	assert.Len(t, changes, 1)
	assert.Equal(t, "wallet1", changes[0].WalletAddress)
	assert.Equal(t, "token1", changes[0].TokenMint)
	assert.Equal(t, uint64(1000), changes[0].OldBalance)
	assert.Equal(t, uint64(2000), changes[0].NewBalance)
	assert.Equal(t, 100.0, changes[0].ChangePercent)
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Positive number",
			input:    5.5,
			expected: 5.5,
		},
		{
			name:     "Negative number",
			input:    -5.5,
			expected: 5.5,
		},
		{
			name:     "Zero",
			input:    0.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTokenAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   uint64
		decimals uint8
		expected string
	}{
		{
			name:     "No decimals",
			amount:   1000,
			decimals: 0,
			expected: "1000",
		},
		{
			name:     "With decimals",
			amount:   1000000000,
			decimals: 9,
			expected: "1.0000",
		},
		{
			name:     "Millions",
			amount:   5000000000000,
			decimals: 9,
			expected: "5.00M",
		},
		{
			name:     "Thousands",
			amount:   5000000000,
			decimals: 9,
			expected: "5.00K",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTokenAmount(tt.amount, tt.decimals)
			assert.Equal(t, tt.expected, result)
		})
	}
}
