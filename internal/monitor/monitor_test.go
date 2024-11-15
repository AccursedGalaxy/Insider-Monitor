package monitor

import (
	"testing"
)

func TestDetectChanges(t *testing.T) {
	tests := []struct {
		name          string
		oldData       map[string]*WalletData
		newData       map[string]*WalletData
		expectedCount int
		expectedTypes []string
	}{
		{
			name: "new wallet detection",
			oldData: map[string]*WalletData{},
			newData: map[string]*WalletData{
				"Wallet1": {
					WalletAddress: "Wallet1",
					TokenAccounts: map[string]TokenAccountInfo{
						"TokenA": {Balance: 100},
					},
				},
			},
			expectedCount: 1,
			expectedTypes: []string{"new_wallet"},
		},
		{
			name: "balance change detection",
			oldData: map[string]*WalletData{
				"Wallet1": {
					WalletAddress: "Wallet1",
					TokenAccounts: map[string]TokenAccountInfo{
						"TokenA": {Balance: 100},
					},
				},
			},
			newData: map[string]*WalletData{
				"Wallet1": {
					WalletAddress: "Wallet1",
					TokenAccounts: map[string]TokenAccountInfo{
						"TokenA": {Balance: 200},
					},
				},
			},
			expectedCount: 1,
			expectedTypes: []string{"balance_change"},
		},
		{
			name: "new token detection",
			oldData: map[string]*WalletData{
				"Wallet1": {
					WalletAddress: "Wallet1",
					TokenAccounts: map[string]TokenAccountInfo{
						"TokenA": {Balance: 100},
					},
				},
			},
			newData: map[string]*WalletData{
				"Wallet1": {
					WalletAddress: "Wallet1",
					TokenAccounts: map[string]TokenAccountInfo{
						"TokenA": {Balance: 100},
						"TokenB": {Balance: 200},
					},
				},
			},
			expectedCount: 1,
			expectedTypes: []string{"new_token"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := DetectChanges(tt.oldData, tt.newData)
			
			if len(changes) != tt.expectedCount {
				t.Errorf("expected %d changes, got %d", tt.expectedCount, len(changes))
			}
			
			for i, expectedType := range tt.expectedTypes {
				if i >= len(changes) {
					t.Errorf("missing expected change type: %s", expectedType)
					continue
				}
				if changes[i].ChangeType != expectedType {
					t.Errorf("expected change type %s, got %s", expectedType, changes[i].ChangeType)
				}
			}
		})
	}
} 