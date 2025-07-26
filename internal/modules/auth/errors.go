package auth

import (
	"errors"
	srverr "protravel-finance/internal/shared/server_error"
)

const (
	ServiceErrorIncorrectLoginOrPassword srverr.ErrorTypeUnauthorized = "incorrect login or password"
)

var (
	RepositoryErrorRefreshTokenNotFound = errors.New("refresh token not found")
)
