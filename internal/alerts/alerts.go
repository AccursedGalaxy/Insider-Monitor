package alerts

import (
	"time"
)

type AlertLevel string

const (
	Info     AlertLevel = "INFO"
	Warning  AlertLevel = "WARNING"
	Critical AlertLevel = "CRITICAL"
)

type Alert struct {
	Timestamp     time.Time
	WalletAddress string
	TokenMint     string
	AlertType     string
	Message       string
	Level         AlertLevel
	Data          map[string]interface{} // Additional data for formatting
}

type Alerter interface {
	SendAlert(alert Alert) error
}

// ConsoleAlerter implementation moved to console.go
