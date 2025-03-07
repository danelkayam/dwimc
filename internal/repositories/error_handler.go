package repositories

import (
	"database/sql"
	"dwimc/internal/model"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func handleSQLError(msg string, err error) error {
	destError := model.ErrOperationFailed

	if errors.Is(err, sql.ErrNoRows) {
		destError = model.ErrItemNotFound
	}

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			destError = model.ErrItemAlreadyExists
		}
	}

	return fmt.Errorf("%s: %w (original error: %v)", msg, destError, err)
}
