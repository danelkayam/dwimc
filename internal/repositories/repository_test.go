package repositories

import (
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	goosePath := os.Getenv("GOOSE_PATH")

	if migrationsDir == "" || goosePath == "" {
		t.Fatalf("The envs: MIGRATIONS_DIR, GOOSE_PATH are required!")
	}

	tmpFile, err := os.CreateTemp("", "dwimc-testdb-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp database file: %v", err)
	}

	dsn := tmpFile.Name()
	t.Logf(`Setting up database...
				goose=%s
				migrations=%s
				database=%s`, goosePath, migrationsDir, dsn)

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	cmd := exec.Command(goosePath, "-dir", migrationsDir, "sqlite3", dsn, "up")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Goose migration up failed: %v", err)

		tmpFile.Close()
	}

	t.Cleanup(func() {
		defer db.Close()
		defer os.Remove(dsn)

		cmd := exec.Command(goosePath, "-dir", migrationsDir, "sqlite3", dsn, "down")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Goose migration down failed: %v", err)
		}
	})

	return db
}
