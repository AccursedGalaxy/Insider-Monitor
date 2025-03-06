package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/accursedgalaxy/insider-monitor/internal/auth"
	"github.com/accursedgalaxy/insider-monitor/internal/config"
	"github.com/accursedgalaxy/insider-monitor/internal/monitor"
	"github.com/accursedgalaxy/insider-monitor/internal/storage"
	"github.com/gorilla/mux"
)

// Embed files from the templates and static directories
//
//go:embed templates/* static/*
var content embed.FS

// Server represents the web UI server
type Server struct {
	config     *config.Config
	monitor    *monitor.WalletMonitor
	storage    *storage.Storage
	router     *mux.Router
	templates  *template.Template
	walletData map[string]*WalletData
	port       int
	auth       *auth.Auth
}

// WalletData is a local copy of monitor.WalletData for use in the web package
type WalletData struct {
	WalletAddress string                      `json:"wallet_address"`
	TokenAccounts map[string]TokenAccountInfo `json:"token_accounts"`
	LastScanned   time.Time                   `json:"last_scanned"`
	NetworkURL    string                      `json:"network_url"`
}

// TokenAccountInfo is a local copy of monitor.TokenAccountInfo for use in the web package
type TokenAccountInfo struct {
	Mint        string    `json:"mint"`
	Balance     uint64    `json:"balance"`
	Decimals    uint8     `json:"decimals"`
	Symbol      string    `json:"symbol"`
	LastUpdated time.Time `json:"last_updated"`
}

// convertMonitorData converts monitor.WalletData to web.WalletData
func convertMonitorData(monitorData map[string]*monitor.WalletData) map[string]*WalletData {
	result := make(map[string]*WalletData)

	for addr, data := range monitorData {
		if data == nil {
			continue
		}

		webData := &WalletData{
			WalletAddress: data.WalletAddress,
			TokenAccounts: make(map[string]TokenAccountInfo),
			LastScanned:   data.LastScanned,
		}

		for mint, tokenInfo := range data.TokenAccounts {
			webData.TokenAccounts[mint] = TokenAccountInfo{
				Mint:        mint,
				Balance:     tokenInfo.Balance,
				Decimals:    tokenInfo.Decimals,
				Symbol:      tokenInfo.Symbol,
				LastUpdated: tokenInfo.LastUpdated,
			}
		}

		result[addr] = webData
	}

	return result
}

// NewServer creates a new web server
func NewServer(cfg *config.Config, monitor *monitor.WalletMonitor, storage *storage.Storage, port int) *Server {
	router := mux.NewRouter()
	authService := auth.New()

	server := &Server{
		config:     cfg,
		monitor:    monitor,
		storage:    storage,
		router:     router,
		walletData: make(map[string]*WalletData),
		port:       port,
		auth:       authService,
	}

	// Initialize routes
	server.routes()

	return server
}

// routes initializes the HTTP routes
func (s *Server) routes() {
	// Serve static files from the embedded filesystem with proper subdir
	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatalf("Failed to create static sub-filesystem: %v", err)
	}

	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Page routes
	s.router.HandleFunc("/", s.handleHome).Methods("GET")
	s.router.HandleFunc("/wallets", s.handleWallets).Methods("GET")
	s.router.HandleFunc("/wallets/{address}", s.handleWalletDetail).Methods("GET")
	s.router.HandleFunc("/config", s.handleConfig).Methods("GET")

	// API routes
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Public API endpoints
	apiRouter.HandleFunc("/wallets", s.apiGetWallets).Methods("GET")
	apiRouter.HandleFunc("/wallets/{address}", s.apiGetWalletDetail).Methods("GET")
	apiRouter.HandleFunc("/refresh", s.apiRefreshData).Methods("POST")

	// Authentication endpoint
	apiRouter.HandleFunc("/login", s.handleLogin).Methods("POST")

	// Protected API endpoints (require authentication)
	protectedAPI := apiRouter.PathPrefix("/admin").Subrouter()
	protectedAPI.Use(s.auth.JWTAuthMiddleware)

	protectedAPI.HandleFunc("/config", s.apiGetConfig).Methods("GET")
	protectedAPI.HandleFunc("/config", s.apiUpdateConfig).Methods("PUT", "PATCH")
	protectedAPI.HandleFunc("/wallets", s.apiAddWallet).Methods("POST")
	protectedAPI.HandleFunc("/wallets/{address}", s.apiDeleteWallet).Methods("DELETE")
}

// Start starts the web server
func (s *Server) Start() error {
	log.Printf("Starting web server on port %d...", s.port)

	// Parse templates with function map
	var err error
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"lastKey": func(m interface{}) string {
			var last string
			switch v := m.(type) {
			case map[string]interface{}:
				for k := range v {
					last = k
				}
			case map[string]config.WalletConfig:
				for k := range v {
					last = k
				}
			}
			return last
		},
		"isLastKey": func(m interface{}, key string) bool {
			last := ""
			switch v := m.(type) {
			case map[string]interface{}:
				for k := range v {
					last = k
				}
			case map[string]config.WalletConfig:
				for k := range v {
					last = k
				}
			}
			return last == key
		},
	}

	// First parse all templates
	templatesFS, err := fs.Sub(content, "templates")
	if err != nil {
		return fmt.Errorf("failed to create templates sub-filesystem: %w", err)
	}

	// Parse all templates at once with the function map
	s.templates, err = template.New("").Funcs(funcMap).ParseFS(templatesFS, "*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Print loaded templates for debugging
	for _, t := range s.templates.Templates() {
		log.Printf("Loaded template: %s", t.Name())
	}

	// Load initial wallet data
	data, err := s.storage.LoadWalletData()
	if err == nil {
		s.walletData = convertMonitorData(data)
	}

	// Start HTTP server
	srv := &http.Server{
		Handler:      s.router,
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

// Handlers for pages

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Dashboard - Solana Insider Monitor",
		"ActivePage": "home",
		"WalletData": s.walletData,
		"Config":     s.config,
	}

	s.renderTemplate(w, "index.html", data)
}

func (s *Server) handleWallets(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Wallets - Solana Insider Monitor",
		"ActivePage": "wallets",
		"WalletData": s.walletData,
		"Config":     s.config,
	}

	s.renderTemplate(w, "wallets.html", data)
}

func (s *Server) handleWalletDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	walletData, ok := s.walletData[address]
	if !ok {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title":      fmt.Sprintf("Wallet %s - Solana Insider Monitor", address),
		"ActivePage": "wallets",
		"Wallet":     walletData,
		"Address":    address,
		"Config":     s.config,
	}

	s.renderTemplate(w, "wallet_detail.html", data)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "Configuration - Solana Insider Monitor",
		"ActivePage": "config",
		"Config":     s.config,
	}

	s.renderTemplate(w, "config.html", data)
}

// API Handlers

func (s *Server) apiGetWallets(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, s.walletData)
}

func (s *Server) apiGetWalletDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	walletData, ok := s.walletData[address]
	if !ok {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	respondJSON(w, walletData)
}

func (s *Server) apiRefreshData(w http.ResponseWriter, r *http.Request) {
	// Run a scan
	newData, err := s.monitor.ScanAllWallets()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan wallets: %v", err), http.StatusInternalServerError)
		return
	}

	// Update wallet data
	s.walletData = convertMonitorData(newData)

	// Save to storage
	if err := s.storage.SaveWalletData(newData); err != nil {
		log.Printf("Warning: Failed to save wallet data: %v", err)
	}

	respondJSON(w, map[string]string{"status": "success", "message": "Data refreshed"})
}

// Handle login and generate JWT token
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := s.auth.GenerateToken(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	respondJSON(w, map[string]string{
		"token": token,
	})
}

// Get current configuration
func (s *Server) apiGetConfig(w http.ResponseWriter, r *http.Request) {
	// Create a response-specific copy of config with fields matching JavaScript expectations
	configResponse := struct {
		NetworkURL   string   `json:"NetworkURL"`
		Wallets      []string `json:"Wallets"`
		ScanInterval string   `json:"ScanInterval"`
		Alerts       struct {
			MinimumBalance    float64  `json:"MinimumBalance"`
			SignificantChange float64  `json:"SignificantChange"`
			IgnoreTokens      []string `json:"IgnoreTokens"`
		} `json:"Alerts"`
		Discord struct {
			Enabled    bool   `json:"Enabled"`
			WebhookURL string `json:"WebhookURL"`
			ChannelID  string `json:"ChannelID"`
		} `json:"Discord"`
	}{
		NetworkURL:   s.config.NetworkURL,
		Wallets:      s.config.Wallets,
		ScanInterval: s.config.ScanInterval,
		Alerts: struct {
			MinimumBalance    float64  `json:"MinimumBalance"`
			SignificantChange float64  `json:"SignificantChange"`
			IgnoreTokens      []string `json:"IgnoreTokens"`
		}{
			MinimumBalance:    s.config.Alerts.MinimumBalance,
			SignificantChange: s.config.Alerts.SignificantChange,
			IgnoreTokens:      s.config.Alerts.IgnoreTokens,
		},
		Discord: struct {
			Enabled    bool   `json:"Enabled"`
			WebhookURL string `json:"WebhookURL"`
			ChannelID  string `json:"ChannelID"`
		}{
			Enabled:    s.config.Discord.Enabled,
			WebhookURL: s.config.Discord.WebhookURL,
			ChannelID:  s.config.Discord.ChannelID,
		},
	}

	respondJSON(w, configResponse)
}

// Update configuration
func (s *Server) apiUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var update config.UpdateRequest

	// Read and log the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the raw request
	log.Printf("Received update request: %s", string(body))

	// Reset the body for further processing
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Log the parsed update
	log.Printf("Parsed update: %+v", update)

	if err := s.config.Update(update); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update config: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the updated config
	log.Printf("Updated config: %+v", s.config)

	// If the scan interval was updated, we need to restart the monitor with the new interval
	if update.ScanInterval != nil {
		// Parse the scan interval string to a Duration
		interval, err := time.ParseDuration(s.config.ScanInterval)
		if err != nil {
			log.Printf("Warning: Invalid scan interval format: %v", err)
			// Either continue with default or return an error
		} else {
			// Pass the duration value to the method
			s.monitor.UpdateScanInterval(interval)
		}
	}

	respondJSON(w, map[string]string{
		"status":  "success",
		"message": "Configuration updated successfully",
	})
}

// Add a new wallet to monitor
func (s *Server) apiAddWallet(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate wallet address (simple validation)
	if len(request.Address) != 44 && len(request.Address) != 43 {
		http.Error(w, "Invalid wallet address", http.StatusBadRequest)
		return
	}

	// Check if wallet already exists
	for _, wallet := range s.config.Wallets {
		if wallet == request.Address {
			http.Error(w, "Wallet already exists", http.StatusBadRequest)
			return
		}
	}

	// Add wallet to config
	s.config.Wallets = append(s.config.Wallets, request.Address)

	// Save config
	if err := s.config.Save(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	// Add wallet to monitor
	if err := s.monitor.AddWallet(request.Address); err != nil {
		// Failed to add to monitor, but config is already updated
		log.Printf("Warning: Failed to add wallet to monitor: %v", err)
	}

	respondJSON(w, map[string]string{
		"status":  "success",
		"message": "Wallet added successfully",
	})
}

// Delete a wallet from monitoring
func (s *Server) apiDeleteWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	// Find wallet in config
	found := false
	var newWallets []string

	for _, wallet := range s.config.Wallets {
		if wallet != address {
			newWallets = append(newWallets, wallet)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Update config
	s.config.Wallets = newWallets

	// Save config
	if err := s.config.Save(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	// Remove wallet from monitor
	s.monitor.RemoveWallet(address)

	respondJSON(w, map[string]string{
		"status":  "success",
		"message": "Wallet removed successfully",
	})
}

// Helper methods

func (s *Server) renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	// Set the content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Debug template rendering
	log.Printf("Rendering template: %s", tmpl)

	// Execute the layout template with the data
	err := s.templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}
