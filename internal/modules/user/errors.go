package user

import (
	"errors"
	srverr "protravel-finance/internal/shared/server_error"
)

// Service Error
const (
	ServiceErrorCurrencyCodeNotAllowed srverr.ErrorTypeNotFound = "currency code not allowed"

	ServiceErrorUserNotFound      srverr.ErrorTypeNotFound = "user not found"
	ServiceErrorUserAlreadyExists srverr.ErrorTypeConflict = "user already exists"

	ServiceErrorLenPassword        srverr.ErrorTypeBadRequest = "password must be at least 8 characters long"
	ServiceErrorPasswordContent    srverr.ErrorTypeBadRequest = "password must contain at least one letter and one digit"
	ServiceErrorPhoneIsNotValid    srverr.ErrorTypeBadRequest = "phone is not valid"
	ServiceErrorEmailIsNotValid    srverr.ErrorTypeBadRequest = "email is not valid"
	ServiceErrorEmailTooLong       srverr.ErrorTypeBadRequest = "email is too long"
	ServiceErrorEmailAndPhoneEmpty srverr.ErrorTypeBadRequest = "email and phone can't be empty"
)

// Repository Error
var (
	RepositoryErrorUserAlreadyExists = errors.New("user already exists")
	RepositoryErrorUserNotFound      = errors.New("user not found")
)
