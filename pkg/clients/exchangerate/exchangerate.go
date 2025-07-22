package exchangerate

import "context"

type ExchangeRate interface {
	GetExchangeRates(ctx context.Context, base string) (*ExchangeRatesResponse, error)
}

type ExchangeRatesResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}