package middleware

import (
	"net/http"
	"protravel-finance/pkg/logger"
)

type Middleware interface {
	PanicRecovery(next http.Handler) http.Handler
}

type middleware struct {
	log logger.Logger
}

func NewMiddleware(log logger.Logger) Middleware {
	return &middleware{
		log: log,
	}
}

func (m *middleware) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				m.log.Errorf("[%s] Panic recovered: %v", r.URL.String(), rec)
				w.WriteHeader(500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
