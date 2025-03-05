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

	server := &Server{
		config:     cfg,
		monitor:    monitor,
		storage:    storage,
		router:     router,
		walletData: make(map[string]*WalletData),
		port:       port,
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
	s.router.HandleFunc("/api/wallets", s.apiGetWallets).Methods("GET")
	s.router.HandleFunc("/api/wallets/{address}", s.apiGetWalletDetail).Methods("GET")
	s.router.HandleFunc("/api/refresh", s.apiRefreshData).Methods("POST")
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
