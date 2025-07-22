package runner

import (
	"github.com/gorilla/mux"
	"protravel-finance/internal/shared/middleware"
)

type Handler interface {
	Init() []Runner

	RouterWithVersion() *mux.Router
}

type Runner interface {
	Run(router *mux.Router, mid middleware.Middleware)
}

func InitHandlers(
	router *mux.Router,
	mid middleware.Middleware,

	handlers ...Handler,
) {

	run(
		mid,
		handlers...,
	)
	router.Use(mid.PanicRecovery)
}

func run(mid middleware.Middleware, handlers ...Handler) {
	for _, h := range handlers {
		for _, r := range h.Init() {
			r.Run(h.RouterWithVersion(), mid)
		}
	}
}
