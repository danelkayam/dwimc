package repositories

import (
	"testing"
)

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSQLUserRepository(db)

	t.Run("create user", func(t *testing.T) {
		email := "testuser@dwimc.awesome"
		password := "hashedpass"
		token := "some-token"
		createdUser, err := repo.Create(email, password, token)

		if err != nil {
			t.Fatalf("Create User failed: %v", err)
		}

		if createdUser.ID == 0 {
			t.Fatalf("Expected a valid ID, got 0")
		}

		if email != createdUser.Email {
			t.Fatalf("Expected email: %s, got %s", email, createdUser.Email)
		}

		if password != createdUser.Password {
			t.Fatalf("Expected email: %s, got %s", password, createdUser.Password)
		}

		t.Logf("Success creating user: %v", createdUser)
	})
}
