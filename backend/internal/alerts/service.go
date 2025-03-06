package alerts

import (
	"log"
	"sync"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/config"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
)

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	InfoLevel    AlertLevel = "INFO"
	WarningLevel AlertLevel = "WARNING"
	ErrorLevel   AlertLevel = "ERROR"
)

// AlertType represents the type of alert
type AlertType string

const (
	BalanceChangeAlert AlertType = "BALANCE_CHANGE"
	NewTokenAlert      AlertType = "NEW_TOKEN"
	TokenRemovedAlert  AlertType = "TOKEN_REMOVED"
	ScanErrorAlert     AlertType = "SCAN_ERROR"
)

// Alert represents a system alert
type Alert struct {
	ID            string                 `json:"id"`
	Timestamp     time.Time              `json:"timestamp"`
	WalletAddress string                 `json:"wallet_address"`
	TokenMint     string                 `json:"token_mint,omitempty"`
	AlertType     AlertType              `json:"alert_type"`
	Message       string                 `json:"message"`
	Level         AlertLevel             `json:"level"`
	Data          map[string]interface{} `json:"data"`
}

// Service manages alerts
type Service struct {
	alerts        []Alert
	mu            sync.RWMutex
	config        *config.Service
	maxAlerts     int
	deliveryChain []AlertDelivery
}

// AlertDelivery is an interface for alert delivery methods
type AlertDelivery interface {
	DeliverAlert(alert Alert) error
}

var (
	alertService     *Service
	alertServiceOnce sync.Once
)

// GetInstance returns the singleton instance of the alerts service
func GetInstance() *Service {
	alertServiceOnce.Do(func() {
		service := &Service{
			alerts:    make([]Alert, 0, 100),
			config:    config.GetInstance(),
			maxAlerts: 1000, // Keep the last 1000 alerts
		}

		// Initialize alert deliveries
		service.deliveryChain = []AlertDelivery{
			&ConsoleDelivery{}, // Always deliver to console
		}

		// Add Discord delivery if enabled in config
		cfg := service.config.GetConfig()
		if cfg.Discord.Enabled && cfg.Discord.WebhookURL != "" {
			service.deliveryChain = append(service.deliveryChain,
				&DiscordDelivery{
					WebhookURL: cfg.Discord.WebhookURL,
					ChannelID:  cfg.Discord.ChannelID,
				},
			)
		}

		// Connect to monitor events
		service.connectToMonitorEvents()

		// Set the global instance
		alertService = service

		log.Println("Alert service initialized")
	})

	return alertService
}

// connectToMonitorEvents subscribes to monitor events
func (s *Service) connectToMonitorEvents() {
	emitter := monitor.GetEventEmitter()

	// Listen for token changes
	emitter.Subscribe(monitor.TokenChange, func(event monitor.Event) {
		// Process token change event
		s.processTokenChangeEvent(event)
	})

	// Listen for wallet scan errors
	emitter.Subscribe(monitor.WalletScanError, func(event monitor.Event) {
		// Process scan error event
		s.processScanErrorEvent(event)
	})

	log.Println("Alert service connected to monitor events")
}

// processTokenChangeEvent processes token change events and generates alerts
func (s *Service) processTokenChangeEvent(event monitor.Event) {
	// Extract data from event
	walletAddress, ok := event.Payload["wallet_address"].(string)
	if !ok {
		return
	}

	tokenMint, ok := event.Payload["token_mint"].(string)
	if !ok {
		return
	}

	// Get alert configuration
	cfg := s.config.GetConfig()

	// Check if this token is in the ignore list
	for _, ignoredToken := range cfg.Alerts.IgnoreTokens {
		if ignoredToken == tokenMint {
			return
		}
	}

	// Process balance change
	oldBalance, oldOk := event.Payload["old_balance"].(uint64)
	newBalance, newOk := event.Payload["new_balance"].(uint64)

	if oldOk && newOk && oldBalance > 0 {
		// Calculate percentage change
		changePercent := float64(newBalance-oldBalance) / float64(oldBalance)

		// Check if change exceeds threshold
		if changePercent >= cfg.Alerts.SignificantChange || changePercent <= -cfg.Alerts.SignificantChange {
			// Create alert
			alert := Alert{
				ID:            generateAlertID(),
				Timestamp:     time.Now(),
				WalletAddress: walletAddress,
				TokenMint:     tokenMint,
				AlertType:     BalanceChangeAlert,
				Message:       "Token balance changed significantly",
				Level:         InfoLevel,
				Data: map[string]interface{}{
					"previous_balance": oldBalance,
					"new_balance":      newBalance,
					"change_percent":   changePercent * 100, // Convert to percentage
					"symbol":           event.Payload["symbol"],
				},
			}

			// Add the alert
			s.addAlert(alert)
		}
	} else if newOk && newBalance > 0 && !oldOk {
		// New token added
		alert := Alert{
			ID:            generateAlertID(),
			Timestamp:     time.Now(),
			WalletAddress: walletAddress,
			TokenMint:     tokenMint,
			AlertType:     NewTokenAlert,
			Message:       "New token detected in wallet",
			Level:         WarningLevel,
			Data: map[string]interface{}{
				"balance": newBalance,
				"symbol":  event.Payload["symbol"],
			},
		}

		// Add the alert
		s.addAlert(alert)
	}
}

// processScanErrorEvent processes scan error events
func (s *Service) processScanErrorEvent(event monitor.Event) {
	errorMsg, ok := event.Payload["error"].(string)
	if !ok {
		return
	}

	// Create alert
	alert := Alert{
		ID:        generateAlertID(),
		Timestamp: time.Now(),
		AlertType: ScanErrorAlert,
		Message:   "Wallet scan failed",
		Level:     ErrorLevel,
		Data: map[string]interface{}{
			"error": errorMsg,
		},
	}

	// Add the alert
	s.addAlert(alert)
}

// addAlert adds an alert to the service and delivers it
func (s *Service) addAlert(alert Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to alerts list
	s.alerts = append(s.alerts, alert)

	// Trim if needed
	if len(s.alerts) > s.maxAlerts {
		s.alerts = s.alerts[len(s.alerts)-s.maxAlerts:]
	}

	// Deliver the alert
	for _, delivery := range s.deliveryChain {
		if err := delivery.DeliverAlert(alert); err != nil {
			log.Printf("Alert delivery failed: %v", err)
		}
	}
}

// GetAlerts returns all stored alerts
func (s *Service) GetAlerts() []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to return
	alertsCopy := make([]Alert, len(s.alerts))
	copy(alertsCopy, s.alerts)

	return alertsCopy
}

// GetAlertsByWallet returns alerts for a specific wallet
func (s *Service) GetAlertsByWallet(walletAddress string) []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var walletAlerts []Alert
	for _, alert := range s.alerts {
		if alert.WalletAddress == walletAddress {
			walletAlerts = append(walletAlerts, alert)
		}
	}

	return walletAlerts
}

// generateAlertID generates a unique ID for an alert
func generateAlertID() string {
	return "alert-" + time.Now().Format("20060102-150405-000")
}
