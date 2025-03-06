package alerts

import (
	"errors"
	"testing"
	"time"
)

// MockAlertDelivery captures alerts for testing
type MockAlertDelivery struct {
	Alerts []Alert
}

// DeliverAlert captures an alert
func (d *MockAlertDelivery) DeliverAlert(alert Alert) error {
	d.Alerts = append(d.Alerts, alert)
	return nil
}

// FailingAlertDelivery always fails to deliver alerts
type FailingAlertDelivery struct{}

// DeliverAlert always returns an error
func (d *FailingAlertDelivery) DeliverAlert(alert Alert) error {
	return errors.New("mock delivery failure")
}

func TestAddAlert(t *testing.T) {
	// Create a service with mock delivery
	mockDelivery := &MockAlertDelivery{
		Alerts: make([]Alert, 0),
	}

	service := &Service{
		alerts:        make([]Alert, 0),
		maxAlerts:     10,
		deliveryChain: []AlertDelivery{mockDelivery},
	}

	// Create a test alert
	alert := Alert{
		ID:            "test-alert",
		Timestamp:     time.Now(),
		WalletAddress: "test-wallet",
		TokenMint:     "test-token",
		AlertType:     BalanceChangeAlert,
		Message:       "Test alert",
		Level:         InfoLevel,
		Data: map[string]interface{}{
			"test": "data",
		},
	}

	// Add the alert
	service.addAlert(alert)

	// Check if the alert was added
	if len(service.alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(service.alerts))
	}

	// Check if the alert was delivered
	if len(mockDelivery.Alerts) != 1 {
		t.Errorf("Expected 1 delivered alert, got %d", len(mockDelivery.Alerts))
	}

	// Check that the alert content was preserved
	if mockDelivery.Alerts[0].ID != "test-alert" {
		t.Errorf("Expected alert ID test-alert, got %s", mockDelivery.Alerts[0].ID)
	}
}

func TestAddAlertWithFailingDelivery(t *testing.T) {
	// Create a service with a failing delivery
	service := &Service{
		alerts:        make([]Alert, 0),
		maxAlerts:     10,
		deliveryChain: []AlertDelivery{&FailingAlertDelivery{}},
	}

	// Create a test alert
	alert := Alert{
		ID:        "test-alert",
		Timestamp: time.Now(),
		AlertType: BalanceChangeAlert,
		Message:   "Test alert",
		Level:     InfoLevel,
	}

	// Add the alert - this should not panic despite delivery failure
	service.addAlert(alert)

	// Check if the alert was still added
	if len(service.alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(service.alerts))
	}
}

func TestMaxAlerts(t *testing.T) {
	// Create a service with a small max alerts limit
	service := &Service{
		alerts:        make([]Alert, 0),
		maxAlerts:     3,
		deliveryChain: []AlertDelivery{&MockAlertDelivery{Alerts: make([]Alert, 0)}},
	}

	// Add more than maxAlerts alerts
	for i := 0; i < 5; i++ {
		alert := Alert{
			ID:        generateAlertID(),
			Timestamp: time.Now(),
			AlertType: BalanceChangeAlert,
			Message:   "Test alert",
			Level:     InfoLevel,
		}
		service.addAlert(alert)
	}

	// Check that we only kept maxAlerts alerts
	if len(service.alerts) != 3 {
		t.Errorf("Expected %d alerts, got %d", service.maxAlerts, len(service.alerts))
	}
}

func TestGetAlerts(t *testing.T) {
	// Create a service with some test alerts
	service := &Service{
		alerts:    make([]Alert, 0),
		maxAlerts: 10,
	}

	// Add test alerts
	for i := 0; i < 3; i++ {
		alert := Alert{
			ID:        generateAlertID(),
			Timestamp: time.Now(),
			AlertType: BalanceChangeAlert,
			Message:   "Test alert",
			Level:     InfoLevel,
		}
		service.alerts = append(service.alerts, alert)
	}

	// Get the alerts
	alerts := service.GetAlerts()

	// Check that we got all alerts
	if len(alerts) != 3 {
		t.Errorf("Expected 3 alerts, got %d", len(alerts))
	}
}

func TestGetAlertsByWallet(t *testing.T) {
	// Create a service with alerts for different wallets
	service := &Service{
		alerts:    make([]Alert, 0),
		maxAlerts: 10,
	}

	// Add alerts for wallet1
	for i := 0; i < 2; i++ {
		alert := Alert{
			ID:            generateAlertID(),
			Timestamp:     time.Now(),
			WalletAddress: "wallet1",
			AlertType:     BalanceChangeAlert,
			Message:       "Test alert",
			Level:         InfoLevel,
		}
		service.alerts = append(service.alerts, alert)
	}

	// Add alerts for wallet2
	for i := 0; i < 1; i++ {
		alert := Alert{
			ID:            generateAlertID(),
			Timestamp:     time.Now(),
			WalletAddress: "wallet2",
			AlertType:     BalanceChangeAlert,
			Message:       "Test alert",
			Level:         InfoLevel,
		}
		service.alerts = append(service.alerts, alert)
	}

	// Get alerts for wallet1
	wallet1Alerts := service.GetAlertsByWallet("wallet1")

	// Check that we got the right number of alerts
	if len(wallet1Alerts) != 2 {
		t.Errorf("Expected 2 alerts for wallet1, got %d", len(wallet1Alerts))
	}

	// Get alerts for wallet2
	wallet2Alerts := service.GetAlertsByWallet("wallet2")

	// Check that we got the right number of alerts
	if len(wallet2Alerts) != 1 {
		t.Errorf("Expected 1 alert for wallet2, got %d", len(wallet2Alerts))
	}

	// Get alerts for non-existent wallet
	nonExistentAlerts := service.GetAlertsByWallet("wallet3")

	// Check that we got no alerts
	if len(nonExistentAlerts) != 0 {
		t.Errorf("Expected 0 alerts for wallet3, got %d", len(nonExistentAlerts))
	}
}
