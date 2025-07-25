package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/shared/database/postgres"
	srverr "protravel-finance/internal/shared/server_error"
	"protravel-finance/pkg/generate"
	"regexp"
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

var (
	passwordRegexp = regexp.MustCompile(`[a-zA-Z].*\d|\d.*[a-zA-Z]`)
)

func (m *userService) CreateUser(ctx context.Context, tx pgx.Tx, user *domain.User, password string) (*domain.User, srverr.ServerError) {
	if !domain.CurrencyCodesMap[user.PreferredCurrency] {
		return nil, srverr.NewServerError(ServiceErrorCurrencyCodeNotAllowed)
	}
	if utf8.RuneCountInString(password) < 8 {
		return nil, srverr.NewServerError(ServiceErrorLenPassword)
	}
	if user.Email == "" && user.Phone == "" {
		return nil, srverr.NewServerError(ServiceErrorEmailAndPhoneEmpty)
	}

	if user.Phone != "" {
		if matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, user.Phone); !matched {
			return nil, srverr.NewServerError(ServiceErrorPhoneIsNotValid)
		}
	}
	if user.Email != "" {
		_, err := mail.ParseAddress(user.Email)
		if err != nil {
			return nil, srverr.NewServerError(ServiceErrorEmailIsNotValid)
		}
		if len(user.Email) > 255 {
			return nil, srverr.NewServerError(ServiceErrorEmailTooLong)
		}
	}
	if matched := passwordRegexp.MatchString(password); !matched {
		return nil, srverr.NewServerError(ServiceErrorPasswordContent)
	}

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
