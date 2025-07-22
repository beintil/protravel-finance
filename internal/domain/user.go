package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                uuid.UUID
	PublicID          string
	Email             string
	PasswordHash      string
	FirstName         string
	LastName          string
	PreferredCurrency string
	Language          string
	Timezone          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
