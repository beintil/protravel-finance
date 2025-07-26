package transperr

import (
	"net/http"
	srverr "protravel-finance/internal/shared/server_error"
)

type ErrorConverter interface {
	ToHTTP(err srverr.ServerError) TransportError
	// ToGRPC, etc..
}

type errorConverter struct{}

func NewErrorConverter() ErrorConverter {
	return &errorConverter{}
}

func (e *errorConverter) ToHTTP(err srverr.ServerError) TransportError {
	var code int

	switch err.GetServerError().(type) {
	case srverr.ErrorTypeConflict:
		code = http.StatusConflict
	case srverr.ErrorTypeBadRequest:
		code = http.StatusBadRequest
	case srverr.ErrorTypeNotFound:
		code = http.StatusNotFound
	case srverr.ErrorTypeUnauthorized:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	tErr := NewTransportError(err.Error(), code)

	if code != http.StatusInternalServerError {
		_ = tErr.SetMessage(err.GetServerError().String()).SetDetails(err.GetDetails())
	}
	return tErr
}
