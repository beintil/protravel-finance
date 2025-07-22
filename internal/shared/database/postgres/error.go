package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresError string

const (
	DuplicateKeyValueViolatesUniqueConstraint PostgresError = "23505"
	ErrNoRows                                 PostgresError = "23503"
)

func ErrorIs(err error, target PostgresError) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && (pgErr.Code == string(target))
}
