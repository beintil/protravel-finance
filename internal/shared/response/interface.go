package response

import (
	"net/http"
	"protravel-finance/models"
)

type HttpResponse interface {
	ErrorResponse(w http.ResponseWriter, r *http.Request, err *models.TransportError)
	WriteResponse(w http.ResponseWriter, r *http.Request, code int, resp any)
}
