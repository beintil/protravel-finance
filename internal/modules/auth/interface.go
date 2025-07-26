package auth

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net/http"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/runner"
	srverr "protravel-finance/internal/shared/server_error"
)

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	runner.Runner
}

type Service interface {
	Register(ctx context.Context, registerUser *domain.RegisterUser) (*domain.User, *domain.AuthToken, srverr.ServerError)
	Login(ctx context.Context, loginUser *domain.LoginUser) (*domain.User, *domain.AuthToken, srverr.ServerError)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthToken, srverr.ServerError)
	Logout(ctx context.Context, refreshToken string) srverr.ServerError
}

type Repository interface {
	SaveRefreshToken(ctx context.Context, redisClient *redis.Client, refreshToken string, jti string) error
	GetRefreshToken(ctx context.Context, redisClient *redis.Client, jti string) (string, error)
	DeleteRefreshToken(ctx context.Context, redisClient *redis.Client, jti string) error
}
