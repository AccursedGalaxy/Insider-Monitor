package models

// Config represents the application configuration
type Config struct {
	NetworkURL    string                  `json:"network_url"`
	Wallets       []string                `json:"wallets"`
	ScanInterval  string                  `json:"scan_interval"`
	Scan          ScanConfig              `json:"scan"`
	WalletConfigs map[string]WalletConfig `json:"wallet_configs,omitempty"`
	Alerts        AlertSettings           `json:"alerts"`
}

// SystemStatus represents the current system status
type SystemStatus struct {
	Status           string `json:"status"`
	Uptime           int    `json:"uptime"`
	Version          string `json:"version"`
	BackendVersion   string `json:"backend_version"`
	FrontendVersion  string `json:"frontend_version"`
	LastScan         string `json:"last_scan"`
	NextScan         string `json:"next_scan"`
	ConnectedClients int    `json:"connected_clients"`
	MemoryUsage      string `json:"memory_usage"`
	CPUUsage         string `json:"cpu_usage"`
}

// ScanStatus represents the status of wallet scanning
type ScanStatus struct {
	IsScanning       bool               `json:"is_scanning"`
	ScanInterval     string             `json:"scan_interval"`
	LastScan         string             `json:"last_scan"`
	NextScan         string             `json:"next_scan"`
	MonitoredWallets int                `json:"monitored_wallets"`
	ScannedTokens    int                `json:"scanned_tokens"`
	ScanHistory      []ScanHistoryEntry `json:"scan_history"`
}

// ScanHistoryEntry represents a single scan history entry
type ScanHistoryEntry struct {
	Timestamp  string `json:"timestamp"`
	Duration   string `json:"duration"`
	Successful bool   `json:"successful"`
}
