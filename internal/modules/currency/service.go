package currency

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/shared/database/postgres"
	srverr "protravel-finance/internal/shared/server_error"
	"protravel-finance/pkg/clients/exchangerate"
	"time"
)

type currencyService struct {
	currencyRepos Repository

	transaction postgres.Transaction
	redisClient *redis.Client
}

func NewService(
	currencyRepos Repository,

	transaction postgres.Transaction,
	redisClient *redis.Client,
) Service {
	return &currencyService{
		currencyRepos: currencyRepos,

		transaction: transaction,
		redisClient: redisClient,
	}
}

func (s *currencyService) UpdateCurrencyRates(ctx context.Context, api exchangerate.ExchangeRate) srverr.ServerError {
	resp, err := api.GetExchangeRates(ctx, domain.USDCode.String())
	if err != nil {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error()).SetDetails("failed to get exchange rates from API")
	}
	if len(resp.Rates) == 0 {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("failed to get exchange rates from API").SetDetails("rates is empty")
	}
	var currencyRates = make([]*domain.Currency, 0, len(domain.CurrencyCodes))

	for _, currency := range domain.CurrencyCodes {
		currencyRates = append(currencyRates, &domain.Currency{
			ID:                 uuid.New(),
			BaseCurrencyCode:   domain.USDCode, // Пока что базовая валюта только доллар, так наверняка и останется
			TargetCurrencyCode: currency,
			Rate:               int64(resp.Rates[string(currency)] * domain.CurrencyFromRate),
			Date:               time.Now().UTC(),
		})
	}
	tx, err := s.transaction.BeginTransaction(ctx)
	if err != nil {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	defer s.transaction.Rollback(ctx, tx)

	err = s.currencyRepos.SaveOrUpdateCurrencyRates(ctx, tx, currencyRates)
	if err != nil {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	err = s.currencyRepos.SaveCurrencyRatesInCache(ctx, s.redisClient, currencyRates)
	if err != nil {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}

	err = s.transaction.Commit(ctx, tx)
	if err != nil {
		return srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	return nil
}
