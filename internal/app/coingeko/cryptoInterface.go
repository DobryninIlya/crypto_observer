package coingecko

import (
	"context"
	"cryptoObserver/internal/app/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CryptoInterface interface {
	GetCryptoPrice(ctx context.Context, id string) (*CryptoPriceResponse, error)
}

// CoinGeckoClient реализует взаимодействие с CoinGecko API
type CoinGeckoClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// CryptoPriceResponse представляет структуру ответа от API CoinGecko
type CryptoPriceResponse struct {
	ID           string        `json:"id"`
	Symbol       string        `json:"symbol"`
	Name         string        `json:"name"`
	CurrentPrice model.Decimal `json:"current_price"`
	LastUpdated  string        `json:"last_updated"`
}

// NewCoinGeckoClient создает новый клиент для CoinGecko API
func NewCoinGeckoClient(apiKey string) *CoinGeckoClient {
	return &CoinGeckoClient{
		baseURL: "https://api.coingecko.com/api/v3",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCryptoPrice получает текущую цену криптовалюты по ее ID
func (c *CoinGeckoClient) GetCryptoPrice(ctx context.Context, id string) (*CryptoPriceResponse, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	if c.apiKey != "" {
		req.Header.Add("x-cg-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Проверяем, была ли отмена контекста
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request canceled: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response []CryptoPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("crypto with id %s not found", id)
	}

	return &response[0], nil
}
