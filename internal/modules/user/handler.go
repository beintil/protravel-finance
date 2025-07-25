package user

import (
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"protravel-finance/internal/shared/middleware"
	"protravel-finance/internal/shared/response"
	transperr "protravel-finance/internal/shared/transport_error"
)

type userHandler struct {
	httpResponse response.HttpResponse
	converter    transperr.ErrorConverter

	userService Service

	validationFormat strfmt.Registry
}

func NewHandler(
	httpResponse response.HttpResponse,
	converter transperr.ErrorConverter,

	userService Service,

	validationFormat strfmt.Registry,
) Handler {
	return &userHandler{
		httpResponse: httpResponse,
		converter:    converter,

		userService: userService,

		validationFormat: validationFormat,
	}
}

func (m *userHandler) Run(router *mux.Router, mid middleware.Middleware) {
	userRouter := router.PathPrefix("/user").Subrouter()
	_ = userRouter
}
