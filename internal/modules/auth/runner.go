package auth

import (
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"protravel-finance/internal/runner"
	"protravel-finance/internal/shared/response"
	transperr "protravel-finance/internal/shared/transport_error"
	"protravel-finance/pkg/logger"
	"time"
)

type handlerV1 struct {
	router *mux.Router

	httpResp response.HttpResponse
	log      logger.Logger

	Handler
}

func NewRunnerHandlerV1(
	router *mux.Router,

	httpResp response.HttpResponse,
	converter transperr.ErrorConverter,

	log logger.Logger,
	validationFormat strfmt.Registry,

	authService Service,

	authTimeout time.Duration,
) runner.Handler {
	return &handlerV1{
		router: router.PathPrefix("/v1").Subrouter(),

		httpResp: httpResp,

		log: log,

		Handler: NewHandler(httpResp, converter, authService, validationFormat, authTimeout),
	}
}

func (m *handlerV1) Init() []runner.Runner {
	return []runner.Runner{
		m.Handler,
	}
}

func (m *handlerV1) RouterWithVersion() *mux.Router {
	return m.router
}
