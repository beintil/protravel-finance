package auth

import (
	"context"
	"github.com/go-redis/redis/v8"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/modules/user"
	"protravel-finance/internal/shared/database/postgres"
	srverr "protravel-finance/internal/shared/server_error"
)

type authService struct {
	userService user.Service

	transaction postgres.Transaction
	redisClient *redis.Client
}

func NewAuthService() Service {
	return &authService{}
}

func (m *authService) Register(ctx context.Context, registerUser *domain.RegisterUser) (*domain.User, srverr.ServerError) {
	tx, err := m.transaction.BeginTransaction(ctx)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	defer m.transaction.Rollback(ctx, tx)

	newUser, servErr := m.userService.CreateUser(ctx, tx, registerUser.ToUser(), registerUser.Password)
	if servErr != nil {
		return nil, servErr
	}

	err = m.transaction.Commit(ctx, tx)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	return newUser, nil
}

func (m *authService) Login(ctx context.Context, loginUser *domain.LoginUser) (*domain.User, srverr.ServerError) {
	return nil, nil
}
