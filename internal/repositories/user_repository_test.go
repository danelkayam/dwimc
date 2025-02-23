package repositories

import (
	"fmt"
	"reflect"
	"testing"

	"dwimc/internal/model"
)

type testUser struct {
	id       uint
	email    string
	password string
	token    string
}

func TestUserRepository(t *testing.T) {
	testUsers := generateTestUser()

	t.Run("create users", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		for _, testUser := range testUsers {
			user, err := repo.Create(testUser.email, testUser.password, testUser.token)

			if err != nil {
				t.Fatalf("Create User failed: %v", err)
			}

			if user.ID == 0 {
				t.Fatalf("Expected a valid ID, got 0")
			}

			if user.Email != testUser.email {
				t.Fatalf("Got: %s, expected %s (user ID: %d)",
					user.Email, testUser.email, testUser.id)
			}

			if user.Password != testUser.password {
				t.Fatalf("Got: %s, expected %s (user ID: %d)",
					user.Password, testUser.password, testUser.id)
			}

			if user.Token.String != testUser.token {
				t.Fatalf("Got: %s, expected %s (user ID: %d)",
					user.Token.String, testUser.token, testUser.id)
			}
		}
	})

	t.Run("get user", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		user, err := repo.Create(testUsers[0].email, testUsers[0].password, testUsers[0].token)
		if err != nil {
			t.Fatalf("Create User failed: %v", err)
		}

		retrieved, err := repo.GetBy(user.ID)
		checkGetUser(t, retrieved, user, err)

		retrieved, err = repo.GetByEmail(user.Email)
		checkGetUser(t, retrieved, user, err)
	})

	t.Run("update user", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		_, err := repo.Create(testUsers[0].email, testUsers[0].password, testUsers[0].token)
		if err != nil {
			t.Fatalf("Create User failed: %v", err)
		}

		// TODO - implement this - update password
		// TODO - implement this - update token
		// TODO - impelemnt this - update both
	})

	t.Run("delete user", func(t *testing.T) {
		t.Parallel()
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		user, err := repo.Create(testUsers[0].email, testUsers[0].password, testUsers[0].token)
		if err != nil {
			t.Fatalf("Create User failed: %v", err)
		}

		retrieved, err := repo.GetBy(user.ID)
		checkGetUser(t, retrieved, user, err)

		err = repo.Delete(user.ID)
		if err != nil {
			t.Fatalf("Delete User failed: %v", err)
		}

		retrieved, err = repo.GetBy(user.ID)
		if err != nil {
			t.Fatalf("Get User failed: %v", err)
		}

		if retrieved != nil {
			t.Fatalf("Got %+v, expected nil", retrieved)
		}
	})
}

func generateTestUser() []testUser {
	const size = 10
	testUsers := make([]testUser, size)

	for i := 0; i < len(testUsers); i++ {
		testUsers[i] = testUser{
			id:       uint(i),
			email:    fmt.Sprintf("user-%d@dwimc.awesome", i),
			password: fmt.Sprintf("secret-password-%d", i),
			token:    fmt.Sprintf("token-%d", i),
		}
	}

	return testUsers
}

func checkGetUser(t *testing.T, retrievedUser *model.User, user *model.User, err error) {
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if !reflect.DeepEqual(retrievedUser, user) {
		t.Fatalf("Got %+v, expected %+v", retrievedUser, user)
	}
}
