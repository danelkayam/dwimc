package repositories_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	goosePath := os.Getenv("GOOSE_PATH")

	require.NotEmpty(t, migrationsDir, "Env: MIGRATIONS_DIR is required")
	require.NotEmpty(t, goosePath, "Env: GOOSE_PATH is required")

	tmpFile, err := os.CreateTemp("", "dwimc-testdb-*.sqlite")
	require.NoErrorf(t, err, "Failed to create temp database file: %v", err)

	dsn := tmpFile.Name()
	t.Logf(`Setting up database...
				goose=%s
				migrations=%s
				database=%s`, goosePath, migrationsDir, dsn)

	db, err := sqlx.Open("sqlite3", dsn)
	require.NoErrorf(t, err, "Failed to open database: %v", err)

	cmd := exec.Command(goosePath, "-dir", migrationsDir, "sqlite3", dsn, "up")
	if err := cmd.Run(); err != nil {
		require.NoErrorf(t, err, "Goose migration up failed: %v", err)

		tmpFile.Close()
	}

	t.Cleanup(func() {
		defer db.Close()
		defer os.Remove(dsn)

		cmd := exec.Command(goosePath, "-dir", migrationsDir, "sqlite3", dsn, "down")
		if err := cmd.Run(); err != nil {
			require.NoErrorf(t, err, "Goose migration down failed: %v", err)
		}
	})

	return db
}
