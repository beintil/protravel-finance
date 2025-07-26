package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                uuid.UUID
	PublicID          string
	Login             string
	Email             string
	Phone             string
	PasswordHash      string
	FirstName         string
	LastName          string
	PreferredCurrency CurrencyCode
	Language          Language
	Timezone          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UserAuthData struct {
	ID       uuid.UUID
	PublicID string
}
