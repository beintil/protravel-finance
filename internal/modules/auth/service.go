package auth

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/modules/user"
	"protravel-finance/internal/shared/database/postgres"
	jwtsec "protravel-finance/internal/shared/jwt_sec"
	srverr "protravel-finance/internal/shared/server_error"
)

type authService struct {
	userService user.Service

	authRepository Repository

	transaction postgres.Transaction
	redisClient *redis.Client

	jwtSec *jwtsec.JWT
}

func NewService(
	userService user.Service,

	authRepository Repository,

	transaction postgres.Transaction,
	redisClient *redis.Client,

	jwtSec *jwtsec.JWT,
) Service {
	return &authService{
		userService: userService,

		authRepository: authRepository,

		transaction: transaction,
		redisClient: redisClient,

		jwtSec: jwtSec,
	}
}

func (m *authService) Register(
	ctx context.Context,
	registerUser *domain.RegisterUser,
) (*domain.User, *domain.AuthToken, srverr.ServerError) {
	tx, err := m.transaction.BeginTransaction(ctx)
	if err != nil {
		return nil, nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	defer m.transaction.Rollback(ctx, tx)

	newUser, servErr := m.userService.CreateUserWithTx(ctx, tx, registerUser.ToUser(), registerUser.Password)
	if servErr != nil {
		return nil, nil, servErr
	}

	accessToken, srvErr := m.jwtSec.GenerateAccessToken(newUser.ID.String())
	if srvErr != nil {
		return nil,  nil, srvErr
	}
	refreshToken, jti, srvErr := m.jwtSec.GenerateRefreshToken(newUser.ID.String())
	if srvErr != nil {
		return nil, nil, srvErr
	}

	err = m.authRepository.SaveRefreshToken(ctx, m.redisClient, refreshToken, jti)
	if err != nil {
		return nil, nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	err = m.transaction.Commit(ctx, tx)
	if err != nil {
		return nil, nil,  srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}

	return newUser,
		&domain.AuthToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
}

// Login return domain.User and access token
func (m *authService) Login(
	ctx context.Context,
	loginUser *domain.LoginUser,
) (*domain.User, *domain.AuthToken, srverr.ServerError) {
	tx, err := m.transaction.BeginTransaction(ctx)
	if err != nil {
		return nil, nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	defer m.transaction.Rollback(ctx, tx)

	userEntity, srvErr := m.userService.GetUserByLoginParamWithTx(ctx, tx, loginUser.Login)
	if srvErr != nil {
		if srvErr.GetServerError() == user.ServiceErrorUserNotFound {
			return nil, nil, srverr.NewServerError(ServiceErrorIncorrectLoginOrPassword)
		}
		return nil, nil, srvErr
	}
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.PasswordHash), []byte(loginUser.Password))
	if err != nil {
		return nil, nil, srverr.NewServerError(ServiceErrorIncorrectLoginOrPassword)
	}

	accessToken, srvErr := m.jwtSec.GenerateAccessToken(userEntity.ID.String())
	if srvErr != nil {
		return nil, nil, srvErr
	}
	refreshToken, jti, srvErr := m.jwtSec.GenerateRefreshToken(userEntity.ID.String())
	if srvErr != nil {
		return nil, nil, srvErr
	}

	err = m.authRepository.SaveRefreshToken(ctx, m.redisClient, refreshToken, jti)
	if err != nil {
		return nil, nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	err = m.transaction.Commit(ctx, tx)
	if err != nil {
		return nil, nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}

	return userEntity,
		&domain.AuthToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
}

func (m *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthToken, srverr.ServerError) {
	userID, jti, srvErr := m.jwtSec.ValidateRefreshToken(refreshToken)
	if srvErr != nil {
		return nil, srvErr
	}
	storedRefresh, err := m.authRepository.GetRefreshToken(ctx, m.redisClient, jti)
	if err != nil {
		if errors.Is(err, RepositoryErrorRefreshTokenNotFound) {
			return nil, srverr.NewServerError(srverr.ErrUnauthorized)
		}
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	if storedRefresh != refreshToken {
		return nil, srverr.NewServerError(srverr.ErrUnauthorized)
	}
	err = m.authRepository.DeleteRefreshToken(ctx, m.redisClient, jti)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}

	newAccessToken, srvErr := m.jwtSec.GenerateAccessToken(userID)
	if srvErr != nil {
		return nil, srvErr
	}
	newRefreshToken, jti, srvErr := m.jwtSec.GenerateRefreshToken(userID)
	if srvErr != nil {
		return nil, srvErr
	}
	// We generate a new refresh token to minimize the risk of unauthorized access by attackers
	err = m.authRepository.SaveRefreshToken(ctx, m.redisClient, newRefreshToken, jti)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}

	return &domain.AuthToken{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
