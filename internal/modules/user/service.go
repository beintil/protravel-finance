package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/shared/database/postgres"
	srverr "protravel-finance/internal/shared/server_error"
	"protravel-finance/pkg/generate"
	"regexp"
	"strings"
	"unicode/utf8"
)

type userService struct {
	userRepos Repository

	transaction postgres.Transaction
}

func NewService(
	userRepos Repository,
	transaction postgres.Transaction,
) Service {
	return &userService{
		userRepos:   userRepos,
		transaction: transaction,
	}
}

func (m *userService) CreateUser(ctx context.Context, user *domain.User, password string) (*domain.User, srverr.ServerError) {
	if !domain.CurrencyCodesMap[domain.CurrencyCode(user.PreferredCurrency)] {
		return nil, srverr.NewServerError(ServiceErrorCurrencyCodeNotAllowed)
	}
	if !strings.Contains(user.Email, "@") {
		return nil, srverr.NewServerError(ServiceErrorEmailIsNotValid)
	}
	if utf8.RuneCountInString(password) < 8 {
		return nil, srverr.NewServerError(ServiceErrorLenPassword).
			SetError("password must be at least 8 characters long")
	}
	if matched, _ := regexp.MatchString(`[a-zA-Z].*\d|\d.*[a-zA-Z]`, password); !matched {
		return nil, srverr.NewServerError(ServiceErrorPasswordContent)
	}

	tx, err := m.transaction.BeginTransaction(ctx)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	defer m.transaction.Rollback(ctx, tx)

	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		publicID, err := generate.PublicID()
		if err != nil {
			return nil, srverr.NewServerError(srverr.ErrInternalServerError).
				SetError(err.Error())
		}
		_, err = m.userRepos.GetUserByPublicID(ctx, tx, publicID)
		if err != nil {
			if errors.Is(err, RepositoryErrorUserNotFound) {
				user.PublicID = publicID
				break
			}
			return nil, srverr.NewServerError(srverr.ErrInternalServerError).
				SetError(err.Error())
		}
		if i == maxAttempts-1 {
			return nil, srverr.NewServerError(srverr.ErrInternalServerError).
				SetError("failed to generate unique publicID after max attempts")
		}
	}
	user.ID = uuid.New()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	user.PasswordHash = string(hash)

	err = m.userRepos.CreateUser(ctx, tx, user)
	if err != nil {
		if errors.Is(err, RepositoryErrorUserAlreadyExists) {
			return nil, srverr.NewServerError(ServiceErrorUserAlreadyExists)
		}
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	err = m.transaction.Commit(ctx, tx)
	if err != nil {
		return nil, srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error())
	}
	return user, nil
}
