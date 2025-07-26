package user

import (
	"context"
	"github.com/jackc/pgx/v5"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/runner"
	srverr "protravel-finance/internal/shared/server_error"
)

type Handler interface {
	runner.Runner
}

type Service interface {
	CreateUserWithTx(ctx context.Context, tx pgx.Tx, user *domain.User, password string) (*domain.User, srverr.ServerError)
	GetUserByLoginParamWithTx(ctx context.Context, tx pgx.Tx, loginParam string) (*domain.User, srverr.ServerError)
}

type Repository interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *domain.User) error
	GetUserByID(ctx context.Context, tx pgx.Tx, id string) (*domain.User, error)
	GetUserByPublicID(ctx context.Context, tx pgx.Tx, publicID string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.User, error)
	GetUserByLoginParam(ctx context.Context, tx pgx.Tx, loginParam string) (*domain.User, error)
}
