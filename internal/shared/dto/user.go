package dto

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"protravel-finance/internal/domain"
	"protravel-finance/models"
	"protravel-finance/pkg/pointer"
	"time"
)

func UserModelToDomain(m *models.User) *domain.User {
	id, _ := uuid.Parse(string(m.ID))

	return &domain.User{
		ID:                id,
		Email:             m.Email.String(),
		FirstName:         *m.FirstName,
		LastName:          *m.LastName,
		PreferredCurrency: *m.PreferredCurrency,
		Language:          *m.Language,
		Timezone:          *m.Timezone,
		CreatedAt:         time.Time(m.CreatedAt),
		UpdatedAt:         time.Time(m.UpdatedAt),
	}
}

func UserDomainToModel(d *domain.User) *models.User {
	return &models.User{
		ID:                strfmt.UUID(d.ID.String()),
		PublicID:          d.PublicID,
		Email:             pointer.P(strfmt.Email(d.Email)),
		FirstName:         pointer.P(d.FirstName),
		LastName:          pointer.P(d.LastName),
		PreferredCurrency: pointer.P(d.PreferredCurrency),
		Language:          pointer.P(d.Language),
		Timezone:          pointer.P(d.Timezone),
		CreatedAt:         strfmt.DateTime(d.CreatedAt),
		UpdatedAt:         strfmt.DateTime(d.UpdatedAt),
	}
}
