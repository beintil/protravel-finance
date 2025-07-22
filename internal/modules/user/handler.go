package user

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"net/http"
	"protravel-finance/internal/shared/dto"
	"protravel-finance/internal/shared/middleware"
	"protravel-finance/internal/shared/response"
	transperr "protravel-finance/internal/shared/transport_error"
	"protravel-finance/models"
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
	userRouter.HandleFunc("/create", m.CreateUser).Methods(http.MethodPost)
}

func (m *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.User

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			transperr.NewTransportError(
				transperr.ValidationError, http.StatusBadRequest),
		))
		return
	}
	err = req.Validate(m.validationFormat)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r,
			dto.TransportErrorToModel(
				transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
			))
		return
	}
	var user = dto.UserModelToDomain(&req)

	user, errSrv := m.userService.CreateUser(r.Context(), user, *req.Password)
	if errSrv != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			m.converter.ToHTTP(errSrv),
		))
		return
	}
	resp := dto.UserDomainToModel(user)

	m.httpResponse.WriteResponse(w, r, http.StatusCreated, resp)
}
