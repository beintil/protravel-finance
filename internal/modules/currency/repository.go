package currency

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"protravel-finance/internal/domain"
	"strings"
	"time"
)

type currencyRepository struct{}

func NewRepository() Repository {
	return &currencyRepository{}
}

func (m *currencyRepository) SaveOrUpdateCurrencyRates(ctx context.Context, tx pgx.Tx, currencyRates []*domain.Currency) error {
	var values []string
	var args []interface{}

	var columnNum = 5

	for i, rate := range currencyRates {
		values = append(values, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)",
			i*columnNum+1, i*columnNum+2, i*columnNum+3, i*columnNum+4, i*columnNum+5))

		args = append(args,
			rate.ID,
			rate.BaseCurrencyCode,
			rate.TargetCurrencyCode,
			rate.Rate,
			rate.Date,
		)
	}

	sql := fmt.Sprintf(`
		INSERT INTO exchange_rates (id, base_currency, target_currency, rate, date)
		VALUES %s
		ON CONFLICT (id) DO UPDATE SET
			base_currency = EXCLUDED.base_currency,
			target_currency = EXCLUDED.target_currency,
			rate = EXCLUDED.rate,
			date = EXCLUDED.date
	`, strings.Join(values, ","))

	_, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SaveOrUpdateCurrencyRates: %w", err)
	}

	return nil
}

func (m *currencyRepository) SaveCurrencyRatesInCache(ctx context.Context, redisClient *redis.Client, currencyRates []*domain.Currency) error {
	if len(currencyRates) == 0 {
		return fmt.Errorf("SaveCurrencyRatesInCache: currencyRates is empty")
	}
	date := currencyRates[0].Date.Format(time.DateOnly)
	hashKey := fmt.Sprintf("rates:%s", date)

	pipe := redisClient.Pipeline()

	for _, rate := range currencyRates {
		// Field: USD:EUR, Value: 860000
		field := fmt.Sprintf("%s:%s", rate.BaseCurrencyCode, rate.TargetCurrencyCode)
		pipe.HSet(ctx, hashKey, field, rate.Rate)
	}

	// Add date to set for tracking
	pipe.SAdd(ctx, "rates:dates", date)

	// Update current rates pointer
	pipe.Set(ctx, "rates:current", hashKey, 0)

	// Clean old rates (keep last 7 days)
	oldDates, err := redisClient.SMembers(ctx, "rates:dates").Result()
	if err != nil {
		return fmt.Errorf("SaveCurrencyRatesInCache/Result: failed to get dates: %w", err)
	}
	for _, oldDate := range oldDates {
		if parsedDate, err := time.Parse(time.DateOnly, oldDate); err == nil {
			if parsedDate.Before(time.Now().AddDate(0, 0, -7)) {
				pipe.Del(ctx, fmt.Sprintf("rates:%s", oldDate))
				pipe.SRem(ctx, "rates:dates", oldDate)
			}
		}
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("SaveCurrencyRatesInCache/Exec: failed to execute pipeline: %w", err)
	}

	return nil
}
