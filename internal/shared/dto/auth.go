package dto

import (
	"protravel-finance/internal/domain"
	"protravel-finance/models"
)

func LoginUserRequestToDomain(m *models.LoginUserRequest) *domain.LoginUser {
	return &domain.LoginUser{
		Login:    *m.Login,
		Password: m.Password.String(),
	}
}

func LoginUserDomainToModel(user *domain.User, authToken *domain.AuthToken) *models.LoginUserResponse {
	return &models.LoginUserResponse{
		User:      UserDomainToModel(user),
		AuthToken: AuthTokenDomainToModel(authToken),
	}
}

func AuthTokenDomainToModel(m *domain.AuthToken) *models.AuthTokenResponse {
	return &models.AuthTokenResponse{
		AccessToken:  &m.AccessToken,
		RefreshToken: &m.RefreshToken,
	}
}

func RegisterUserRequestToDomain(m *models.RegisterUserRequest) *domain.RegisterUser {
	return &domain.RegisterUser{
		Email:             m.Email.String(),
		Phone:             m.Phone,
		Login:             *m.Login,
		Password:          m.Password.String(),
		PreferredCurrency: domain.CurrencyCode(*m.PreferredCurrency),
		FirstName:         *m.FirstName,
		LastName:          *m.LastName,
		Language:          domain.Language(*m.Language),
		Timezone:          *m.Timezone,
	}
}

func RegisterUserDomainToModel(user *domain.User, token *domain.AuthToken) *models.RegisterUserResponse {
	return &models.RegisterUserResponse{
		AuthToken: AuthTokenDomainToModel(token),
		User:         UserDomainToModel(user),
	}
}
