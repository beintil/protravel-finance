package auth

import (
	"context"
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"net/http"
	"protravel-finance/internal/shared/dto"
	"protravel-finance/internal/shared/middleware"
	"protravel-finance/internal/shared/response"
	transperr "protravel-finance/internal/shared/transport_error"
	"protravel-finance/models"
	"time"
)

type authHandler struct {
	httpResponse response.HttpResponse
	converter    transperr.ErrorConverter

	authService Service

	validationFormat strfmt.Registry

	authTimeout time.Duration
}

func NewHandler(
	httpResponse response.HttpResponse,
	converter transperr.ErrorConverter,

	authService Service,

	validationFormat strfmt.Registry,

	authTimeout time.Duration,
) Handler {
	return &authHandler{
		httpResponse: httpResponse,
		converter:    converter,

		authService: authService,

		validationFormat: validationFormat,

		authTimeout: authTimeout,
	}
}

func (m *authHandler) Run(router *mux.Router, mid middleware.Middleware) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	authRouter.HandleFunc("/login", m.Login).Methods(http.MethodPost)
	authRouter.HandleFunc("/register", m.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", m.RefreshToken).Methods(http.MethodPatch)
	authRouter.HandleFunc("/logout", m.Logout).Methods(http.MethodDelete)
}

func (m *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r,
			dto.TransportErrorToModel(
				transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
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
	ctx, cancel := context.WithTimeout(r.Context(), m.authTimeout)
	defer cancel()

	var loginUser = dto.LoginUserRequestToDomain(&req)

	user, token, srvErr := m.authService.Login(ctx, loginUser)
	if srvErr != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			m.converter.ToHTTP(srvErr),
		))
		return
	}
	var resp = dto.LoginUserDomainToModel(user, token)

	m.httpResponse.WriteResponse(w, r, http.StatusOK, resp)
}

func (m *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r,
			dto.TransportErrorToModel(
				transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
			))
		return
	}
	err = req.Validate(m.validationFormat)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
		))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), m.authTimeout)
	defer cancel()

	var registerUser = dto.RegisterUserRequestToDomain(&req)

	user, token, srvErr := m.authService.Register(ctx, registerUser)
	if srvErr != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			m.converter.ToHTTP(srvErr),
		))
		return
	}
	var resp = dto.RegisterUserDomainToModel(user, token)

	m.httpResponse.WriteResponse(w, r, http.StatusCreated, resp)
}

func (m *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r,
			dto.TransportErrorToModel(
				transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
			))
		return
	}
	err = req.Validate(m.validationFormat)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
		))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), m.authTimeout)
	defer cancel()

	var refreshToken = *req.RefreshToken

	token, srvErr := m.authService.RefreshToken(ctx, refreshToken)
	if srvErr != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			m.converter.ToHTTP(srvErr),
		))
		return
	}
	var resp = dto.AuthTokenDomainToModel(token)

	m.httpResponse.WriteResponse(w, r, http.StatusOK, resp)
}

func (m *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r,
			dto.TransportErrorToModel(
				transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
			))
		return
	}
	err = req.Validate(m.validationFormat)
	if err != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			transperr.NewTransportError(transperr.ValidationError, http.StatusBadRequest),
		))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), m.authTimeout)
	defer cancel()

	var refreshToken = *req.RefreshToken

	srvErr := m.authService.Logout(ctx, refreshToken)
	if srvErr != nil {
		m.httpResponse.ErrorResponse(w, r, dto.TransportErrorToModel(
			m.converter.ToHTTP(srvErr),
		))
		return
	}
	m.httpResponse.WriteResponse(w, r, http.StatusOK, nil)
}
