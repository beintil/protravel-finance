package dto

import (
	"github.com/go-openapi/strfmt"
	"protravel-finance/internal/domain"
	"protravel-finance/models"
)

func UserModelToDomain(m *models.User) *domain.User {
	return &domain.User{
		PublicID:          m.PublicID,
		Email:             m.Email.String(),
		FirstName:         m.FirstName,
		LastName:          m.LastName,
		PreferredCurrency: m.PreferredCurrency,
		Language:          domain.Language(m.Language),
		Timezone:          m.Timezone,
	}
}

func UserDomainToModel(d *domain.User) *models.User {
	return &models.User{
		PublicID:          d.PublicID,
		Email:             strfmt.Email(d.Email),
		FirstName:         d.FirstName,
		LastName:          d.LastName,
		PreferredCurrency: d.PreferredCurrency,
		Language:          d.Language.String(),
		Timezone:          d.Timezone,
	}
}
