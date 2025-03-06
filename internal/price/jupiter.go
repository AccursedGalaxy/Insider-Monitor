package price

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	jupiterPriceV2URL = "https://api.jup.ag/price/v2?ids=%s"
	maxTokensPerBatch = 100 // Jupiter API limit
)

type JupiterPrice struct {
	data  map[string]PriceData
	mutex sync.RWMutex
}

type PriceData struct {
	Price           float64   `json:"price,string"`
	LastUpdated     time.Time `json:"last_updated"`
	ConfidenceLevel string    `json:"confidence_level"`
}

type jupiterResponse struct {
	Data map[string]*struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Price     string `json:"price"`
		ExtraInfo *struct {
			ConfidenceLevel string `json:"confidenceLevel"`
		} `json:"extraInfo,omitempty"`
	} `json:"data"`
	TimeTaken float64 `json:"timeTaken"`
}

func NewJupiterPrice() *JupiterPrice {
	return &JupiterPrice{
		data: make(map[string]PriceData),
	}
}

func (j *JupiterPrice) UpdatePrices(mints []string) error {
	// Split mints into batches of 100 (Jupiter's limit)
	for i := 0; i < len(mints); i += maxTokensPerBatch {
		end := i + maxTokensPerBatch
		if end > len(mints) {
			end = len(mints)
		}

		batch := mints[i:end]
		if err := j.updateBatch(batch); err != nil {
			return fmt.Errorf("failed to update batch %d-%d: %w", i, end, err)
		}

		// Small delay between batches to respect rate limits
		if end < len(mints) {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func (j *JupiterPrice) updateBatch(mints []string) error {
	mintsStr := strings.Join(mints, ",")
	url := fmt.Sprintf(jupiterPriceV2URL, mintsStr)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jupResp jupiterResponse
	if err := json.NewDecoder(resp.Body).Decode(&jupResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	now := time.Now()
	for mint, data := range jupResp.Data {
		if data == nil || data.Price == "" {
			continue
		}

		price, err := parsePrice(data.Price)
		if err != nil {
			continue
		}

		confidence := "medium"
		if data.ExtraInfo != nil {
			confidence = data.ExtraInfo.ConfidenceLevel
		}

		j.data[mint] = PriceData{
			Price:           price,
			LastUpdated:     now,
			ConfidenceLevel: confidence,
		}
	}

	return nil
}

func (j *JupiterPrice) GetPrice(mint string) (PriceData, bool) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()
	data, exists := j.data[mint]
	return data, exists
}

func parsePrice(price string) (float64, error) {
	var value float64
	if _, err := fmt.Sscanf(price, "%f", &value); err != nil {
		return 0, err
	}
	return value, nil
}
