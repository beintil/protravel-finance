package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"protravel-finance/internal/domain"
	"protravel-finance/internal/shared/database/postgres"
)

type userRepository struct{}

func NewRepository() Repository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, tx pgx.Tx, user *domain.User) error {
	_, err := tx.Exec(ctx, `
	INSERT INTO users (
	        id, public_id, email, phone, login,
	    	password_hash, first_name,
	        last_name, preferred_currency,
	        language, timezone
	        ) VALUES (
	        $1, $2, $3, $4, $5,
	        $6, $7, $8, $9, $10, $11 
	    )`,
		user.ID,
		user.PublicID,
		user.Email,
		user.Phone,
		user.Login,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.PreferredCurrency,
		user.Language,
		user.Timezone,
	)
	if err != nil {
		if postgres.ErrorIs(err, postgres.DuplicateKeyValueViolatesUniqueConstraint) {
			return RepositoryErrorUserAlreadyExists
		}
		return fmt.Errorf("CreateUser/Scan: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, tx pgx.Tx, id string) (*domain.User, error) {
	row := tx.QueryRow(ctx, `
	SELECT id, public_id, email, phone, login, password_hash, first_name, last_name, preferred_currency, language, timezone
		FROM users
	WHERE id = $1`,
		id)
	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.Phone,
		&user.Login,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PreferredCurrency,
		&user.Language,
		&user.Timezone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, RepositoryErrorUserNotFound
		}
		return nil, fmt.Errorf("GetUserByID/Scan: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByPublicID(ctx context.Context, tx pgx.Tx, publicID string) (*domain.User, error) {
	row := tx.QueryRow(ctx, `
	SELECT id, public_id, email, phone, login, password_hash, first_name, last_name, preferred_currency, language, timezone
		FROM users
	WHERE public_id = $1`,
		publicID)
	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.Phone,
		&user.Login,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PreferredCurrency,
		&user.Language,
		&user.Timezone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, RepositoryErrorUserNotFound
		}
		return nil, fmt.Errorf("GetUserByPublicID/Scan: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.User, error) {
	row := tx.QueryRow(ctx, `
	SELECT id, public_id, email, phone, login, password_hash, first_name, last_name, preferred_currency, language, timezone
		FROM users
	WHERE email = $1`,
		email)
	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.Phone,
		&user.Login,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PreferredCurrency,
		&user.Language,
		&user.Timezone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, RepositoryErrorUserNotFound
		}
		return nil, fmt.Errorf("GetUserByEmail/Scan: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetUserByLoginParam(ctx context.Context, tx pgx.Tx, loginParam string) (*domain.User, error) {
	var user domain.User

	row := tx.QueryRow(ctx, `
	SELECT id, public_id, email, phone, login, password_hash, first_name, last_name, preferred_currency, language, timezone
		FROM users
	WHERE public_id = $1 OR email = $1 OR phone = $1 OR login = $1`, loginParam)
	err := row.Scan(&user.ID, &user.PublicID, &user.Email, &user.Phone, &user.Login, &user.PasswordHash, &user.FirstName, &user.LastName, &user.PreferredCurrency, &user.Language, &user.Timezone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, RepositoryErrorUserNotFound
		}
		return nil, fmt.Errorf("GetUserByLoginParam/Scan: %w", err)
	}
	return &user, nil
}
