package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	// CurrencyFromRate Множитель для хранения курсов в бд
	CurrencyFromRate = 1e6
)

type Currency struct {
	ID                 uuid.UUID
	BaseCurrencyCode   CurrencyCode // BaseCurrencyCode Базовая валюта [валюта, от которой идут расчёты]
	TargetCurrencyCode CurrencyCode // TargetCurrencyCode Валюта, к которой привязан курс
	Rate               int64        // Rate Курс валюты к USD * CurrencyFromRate
	Date               time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CurrencyCode string

const (
	USDCode CurrencyCode = "USD"
	EURCode CurrencyCode = "EUR"
	RUBCode CurrencyCode = "RUB"
)

var CurrencyCodes = []CurrencyCode{
	USDCode,
	EURCode,
	RUBCode,
}

var CurrencyCodesMap = map[CurrencyCode]bool{
	USDCode: true,
	EURCode: true,
	RUBCode: true,
}
