package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func RecordNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
