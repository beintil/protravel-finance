package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"protravel-finance/pkg/logger"
)

type exchangeRate struct {
	baseURL string
	logger  logger.Logger
	httpClient *http.Client
}

func NewExchangeRate(baseURL string, logger logger.Logger) ExchangeRate {
	return &exchangeRate{
		baseURL: baseURL,
		logger:  logger,
		httpClient: &http.Client{},
	}
}

func (c *exchangeRate) GetExchangeRates(ctx context.Context, base string) (*ExchangeRatesResponse, error) {
	url := fmt.Sprintf("%s/latest/%s", c.baseURL, base)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetExchangeRates/NewRequestWithContext: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetExchangeRates/Do: failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data ExchangeRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("GetExchangeRates/Decode: failed to decode exchange rates: %v", err)
	}

	return &data, nil
}