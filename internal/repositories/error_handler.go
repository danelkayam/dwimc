package repositories

import (
	"database/sql"
	"dwimc/internal/model"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func handleSQLError(err error) error {
	destError := model.ErrOperationFailed

	if errors.Is(err, sql.ErrNoRows) {
		destError = model.ErrItemNotFound
	}

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			destError = model.ErrItemConflict
		}
	}

	return fmt.Errorf("%w (original error: %v)", destError, err)
}
