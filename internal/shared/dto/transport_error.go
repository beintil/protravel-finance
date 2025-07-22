package dto

import (
	"github.com/go-openapi/strfmt"
	transperr "protravel-finance/internal/shared/transport_error"
	"protravel-finance/models"
	"protravel-finance/pkg/pointer"
)

func TransportErrorToModel(m transperr.TransportError) *models.TransportError {
	return &models.TransportError{
		Error:         m.Error(),
		Code:          pointer.P(int32(m.GetCode())),
		Message:       pointer.P(m.GetMessage()),
		Details:       m.GetDetails(),
		TransactionID: strfmt.UUID(m.GetTransactionID().String()),
	}
}
