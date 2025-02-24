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

	t.Run("create user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		t.Run("single valid user", func(t *testing.T) {
			testUser := &testUser{
				email:    "valid-user@dwimc.awesome",
				password: "valid-user-password",
				token:    "valid-user-token",
			}
			user, err := repo.Create(testUser.email, testUser.password, testUser.token)
			checkCreateUser(t, user, testUser, err)
		})

		t.Run("multiple valid users", func(t *testing.T) {
			for _, testUser := range testUsers {
				user, err := repo.Create(testUser.email, testUser.password, testUser.token)
				checkCreateUser(t, user, &testUser, err)
			}
		})

		t.Run("duplicate email", func(t *testing.T) {
			testUser := &testUser{
				email:    "duplicate-email@dwimc.awesome",
				password: "duplicate-email-password",
				token:    "duplicate-email-token",
			}
			repo.Create(testUser.email, testUser.password, testUser.token)
			_, err := repo.Create(testUser.email, "different-password-0", "different-token-0")
			if err == nil {
				t.Fatalf("Got: nil, expected: error for duplicate email")
			}
		})

		t.Run("duplicate token", func(t *testing.T) {
			testUser := &testUser{
				email:    "duplicate-token@dwimc.awesome",
				password: "duplicate-token-password",
				token:    "duplicate-token-token",
			}
			repo.Create(testUser.email, testUser.password, testUser.token)
			_, err := repo.Create("duplicate-token-1@dwimc.awesome", "different-password-1", testUser.token)
			if err == nil {
				t.Fatalf("Got: nil, expected: error for duplicate token")
			}
		})
	})

	t.Run("get user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		t.Run("By ID", func(t *testing.T) {
			t.Parallel()
			user, err := repo.Create(testUsers[0].email, testUsers[0].password, testUsers[0].token)
			if err != nil {
				t.Fatalf("Create User failed: %v", err)
			}

			retrieved, err := repo.GetBy(user.ID)
			checkGetUser(t, retrieved, user, err)
		})

		t.Run("By Email", func(t *testing.T) {
			t.Parallel()
			user, err := repo.Create(testUsers[1].email, testUsers[1].password, testUsers[1].token)
			if err != nil {
				t.Fatalf("Create User failed: %v", err)
			}

			retrieved, err := repo.GetByEmail(user.Email)
			checkGetUser(t, retrieved, user, err)
		})
	})

	t.Run("update user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLUserRepository(db)

		t.Run("update all fields", func(t *testing.T) {
			t.Parallel()

			_, err := repo.Create(testUsers[0].email, testUsers[0].password, testUsers[0].token)
			if err != nil {
				t.Fatalf("Create User failed: %v", err)
			}
		})

		// TODO - implement this - update none
		// TODO - implement this - update password
		// TODO - implement this - update token
		// TODO - implement this - update all
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
			t.Fatalf("Got %+v, expected: nil", retrieved)
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

func checkCreateUser(t *testing.T, createdUser *model.User, testUser *testUser, err error) {
	if err != nil {
		t.Fatalf("Create User failed: %v", err)
	}

	if createdUser.ID == 0 {
		t.Fatalf("Got 0, expected: a valid ID")
	}

	if createdUser.Email != testUser.email {
		t.Fatalf("Got: %s, expected: %s (user ID: %d)",
			createdUser.Email, testUser.email, testUser.id)
	}

	if createdUser.Password != testUser.password {
		t.Fatalf("Got: %s, expected: %s (user ID: %d)",
			createdUser.Password, testUser.password, testUser.id)
	}

	if createdUser.Token.String != testUser.token {
		t.Fatalf("Got: %s, expected: %s (user ID: %d)",
			createdUser.Token.String, testUser.token, testUser.id)
	}
}

func checkGetUser(t *testing.T, retrievedUser *model.User, user *model.User, err error) {
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if !reflect.DeepEqual(retrievedUser, user) {
		t.Fatalf("Got %+v, expected: %+v (user ID: %d)", retrievedUser, user, user.ID)
	}
}
