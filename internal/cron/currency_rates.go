package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"protravel-finance/internal/modules/currency"
	"protravel-finance/pkg/clients/exchangerate"
	"protravel-finance/pkg/logger"
	"time"
)

type updateRatesCron struct {
	currencyService currency.Service

	api exchangerate.ExchangeRate

	logger logger.Logger
}

func NewUpdateRatesCron(
	currencyService currency.Service,
	api exchangerate.ExchangeRate,
	logger logger.Logger,
) Cron {
	return &updateRatesCron{
		currencyService: currencyService,
		logger:          logger,
		api:             api,
	}
}

func (m *updateRatesCron) Run() {
	err := m.currencyService.UpdateCurrencyRates(context.Background(), m.api)
	if err != nil {
		m.logger.Errorf("Failed to update currency rates: %v", err)
		return
	}
	c := cron.New()

	c.AddFunc("0 0 * * * ", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		m.logger.Infof("Start updating currency rates at %s \n", time.Now().UTC().Format(time.DateTime))

		err := m.currencyService.UpdateCurrencyRates(ctx, m.api)
		if err != nil {
			m.logger.Errorf("Failed to update currency rates: %v", err)
			return
		}
		m.logger.Infof("Finished updating currency rates at %s \n", time.Now().UTC().Format(time.DateTime))
	})
	c.Start()
}
