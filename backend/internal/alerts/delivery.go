package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ConsoleDelivery outputs alerts to the console
type ConsoleDelivery struct{}

// DeliverAlert logs the alert to the console
func (d *ConsoleDelivery) DeliverAlert(alert Alert) error {
	log.Printf("[ALERT] %s - %s: %s",
		alert.Level,
		alert.AlertType,
		alert.Message,
	)

	return nil
}

// DiscordDelivery sends alerts to Discord via webhook
type DiscordDelivery struct {
	WebhookURL string
	ChannelID  string
	client     *http.Client
}

// DiscordMessage is the payload for a Discord webhook
type DiscordMessage struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed is a Discord message embed
type DiscordEmbed struct {
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Color       int                    `json:"color,omitempty"` // Decimal value of color
	Timestamp   string                 `json:"timestamp,omitempty"`
	Footer      *DiscordEmbedFooter    `json:"footer,omitempty"`
	Thumbnail   *DiscordEmbedThumbnail `json:"thumbnail,omitempty"`
	Fields      []DiscordEmbedField    `json:"fields,omitempty"`
}

// DiscordEmbedFooter is the footer for a Discord embed
type DiscordEmbedFooter struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// DiscordEmbedThumbnail is a thumbnail for a Discord embed
type DiscordEmbedThumbnail struct {
	URL string `json:"url,omitempty"`
}

// DiscordEmbedField is a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// DeliverAlert sends the alert to Discord
func (d *DiscordDelivery) DeliverAlert(alert Alert) error {
	// Ensure we have a client
	if d.client == nil {
		d.client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	// Create embed fields based on alert data
	var fields []DiscordEmbedField

	// Add wallet address if present
	if alert.WalletAddress != "" {
		fields = append(fields, DiscordEmbedField{
			Name:   "Wallet",
			Value:  fmt.Sprintf("[%s](https://solscan.io/account/%s)", truncateAddress(alert.WalletAddress), alert.WalletAddress),
			Inline: true,
		})
	}

	// Add token mint if present
	if alert.TokenMint != "" {
		fields = append(fields, DiscordEmbedField{
			Name:   "Token",
			Value:  fmt.Sprintf("[%s](https://solscan.io/token/%s)", truncateAddress(alert.TokenMint), alert.TokenMint),
			Inline: true,
		})

		// Add symbol if present
		if symbol, ok := alert.Data["symbol"].(string); ok {
			fields = append(fields, DiscordEmbedField{
				Name:   "Symbol",
				Value:  symbol,
				Inline: true,
			})
		}
	}

	// Add data fields based on alert type
	switch alert.AlertType {
	case BalanceChangeAlert:
		// Add balance change information
		if prevBalance, ok := alert.Data["previous_balance"].(uint64); ok {
			fields = append(fields, DiscordEmbedField{
				Name:   "Previous Balance",
				Value:  fmt.Sprintf("%d", prevBalance),
				Inline: true,
			})
		}

		if newBalance, ok := alert.Data["new_balance"].(uint64); ok {
			fields = append(fields, DiscordEmbedField{
				Name:   "New Balance",
				Value:  fmt.Sprintf("%d", newBalance),
				Inline: true,
			})
		}

		if changePercent, ok := alert.Data["change_percent"].(float64); ok {
			fields = append(fields, DiscordEmbedField{
				Name:   "Change %",
				Value:  fmt.Sprintf("%.2f%%", changePercent),
				Inline: true,
			})
		}

	case ScanErrorAlert:
		// Add error information
		if errorMsg, ok := alert.Data["error"].(string); ok {
			fields = append(fields, DiscordEmbedField{
				Name:   "Error",
				Value:  errorMsg,
				Inline: false,
			})
		}
	}

	// Determine color based on alert level
	var color int
	switch alert.Level {
	case InfoLevel:
		color = 3447003 // Blue
	case WarningLevel:
		color = 16763904 // Orange
	case ErrorLevel:
		color = 15158332 // Red
	default:
		color = 9807270 // Grey
	}

	// Create the message
	message := DiscordMessage{
		Username: "Insider Monitor",
		Embeds: []DiscordEmbed{
			{
				Title:     fmt.Sprintf("%s Alert: %s", alert.AlertType, alert.Message),
				Color:     color,
				Timestamp: alert.Timestamp.Format(time.RFC3339),
				Fields:    fields,
				Footer: &DiscordEmbedFooter{
					Text: "Solana Insider Monitor",
				},
			},
		},
	}

	// Marshal to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord message: %w", err)
	}

	// Send the request
	resp, err := d.client.Post(d.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send Discord message: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord API returned status code %d", resp.StatusCode)
	}

	return nil
}

// truncateAddress shortens an address for display
func truncateAddress(address string) string {
	if len(address) <= 12 {
		return address
	}
	return address[:6] + "..." + address[len(address)-4:]
}
