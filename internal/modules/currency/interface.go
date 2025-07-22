package currency

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"protravel-finance/internal/domain"
	srverr "protravel-finance/internal/shared/server_error"
	"protravel-finance/pkg/clients/exchangerate"
)

type Handler interface{}

type Service interface {
	UpdateCurrencyRates(ctx context.Context, api exchangerate.ExchangeRate) srverr.ServerError
}

type Repository interface {
	SaveOrUpdateCurrencyRates(ctx context.Context, tx pgx.Tx, currencyRates []*domain.Currency) error

	SaveCurrencyRatesInCache(ctx context.Context, redisClient *redis.Client, currencyRates []*domain.Currency) error
}
