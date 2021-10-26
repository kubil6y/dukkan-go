package data

import (
	"errors"

	"github.com/jackc/pgconn"
)

func IsDuplicateRecord(err error) bool {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code == "23505"
		}
	}
	return false
}
