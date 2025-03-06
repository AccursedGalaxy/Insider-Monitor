package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// PublicKey is an alias for solana.PublicKey
type PublicKey = solana.PublicKey

// TokenProgramID is the ID of the token program
var TokenProgramID = solana.TokenProgramID

// PublicKeyFromBase58 converts a base58 string to a PublicKey
func PublicKeyFromBase58(address string) (PublicKey, error) {
	return solana.PublicKeyFromBase58(address)
}

// IsValidAddress checks if a string is a valid Solana address
func IsValidAddress(address string) bool {
	_, err := solana.PublicKeyFromBase58(address)
	return err == nil
}

// FormatTokenAmount formats a token amount with the correct decimals
func FormatTokenAmount(amount uint64, decimals uint8) string {
	if decimals == 0 {
		return fmt.Sprintf("%d", amount)
	}

	divisor := uint64(1)
	for i := uint8(0); i < decimals; i++ {
		divisor *= 10
	}

	whole := amount / divisor
	fraction := amount % divisor

	return fmt.Sprintf("%d.%0*d", whole, int(decimals), fraction)
}

// TokenToDecimal converts a token amount with decimals to a float64 representation
func TokenToDecimal(amount uint64, decimals uint8) float64 {
	if decimals == 0 {
		return float64(amount)
	}

	divisor := float64(1)
	for i := uint8(0); i < decimals; i++ {
		divisor *= 10
	}

	return float64(amount) / divisor
}
