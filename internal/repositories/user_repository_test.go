package repositories_test

import (
	"fmt"
	"testing"
	"time"

	"dwimc/internal/model"
	"dwimc/internal/repositories"
	testutils "dwimc/internal/test_utils"

	"github.com/stretchr/testify/assert"
)

type testUser struct {
	id       uint
	email    string
	password string
}

func TestUserRepository(t *testing.T) {
	testUsers := generateTestUsers()

	t.Run("create user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := repositories.NewSQLUserRepository(db)

		t.Run("single valid user", func(t *testing.T) {
			testUser := &testUser{
				email:    "valid-user@dwimc.awesome",
				password: "valid-user-password",
			}

			user, err := repo.Create(testUser.email, testUser.password)
			assertCreateUser(t, testUser, user, err)
		})

		t.Run("multiple valid users", func(t *testing.T) {
			for _, testUser := range testUsers {
				user, err := repo.Create(testUser.email, testUser.password)
				assertCreateUser(t, &testUser, user, err)
			}
		})

		t.Run("duplicate email", func(t *testing.T) {
			testUser := &testUser{
				email:    "duplicate-email@dwimc.awesome",
				password: "duplicate-email-password",
			}

			_, err := repo.Create(testUser.email, testUser.password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			_, err = repo.Create(testUser.email, "different-password-0")
			assert.ErrorIsf(t, err, model.ErrItemConflict, "Expected error for duplicate email")
		})
	})

	t.Run("get user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := repositories.NewSQLUserRepository(db)

		t.Run("by ID", func(t *testing.T) {
			t.Parallel()

			user, err := repo.Create(testUsers[0].email, testUsers[0].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			retrieved, err := repo.GetByID(user.ID)
			assertGetUser(t, user, retrieved, err)
		})

		t.Run("by ID - not found", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetByID(model.ID(123456789))
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error not found")
		})

		t.Run("by email", func(t *testing.T) {
			t.Parallel()

			user, err := repo.Create(testUsers[1].email, testUsers[1].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			retrieved, err := repo.GetByEmail(user.Email)
			assertGetUser(t, user, retrieved, err)
		})

		t.Run("by email - not found", func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetByEmail("not-existing-user@dwimc.awesome")
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error not found")
		})
	})

	t.Run("update user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := repositories.NewSQLUserRepository(db)

		const UPDATE_SLEEP_DURATION = 1 * time.Second

		t.Run("all fields", func(t *testing.T) {
			t.Parallel()

			createdUser, err := repo.Create(testUsers[0].email, testUsers[0].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			updatedUser, err := repo.Update(createdUser.ID,
				model.WithEmail("updated-email-0"),
				model.WithPassword("updated-password-0"),
			)
			assert.NoErrorf(t, err, "Update User failed: %v", err)

			testutils.AssertEqualItems(createdUser, updatedUser,
				func(field string, shouldBeEqual bool, got any, expected any) {
					t.Helper()
					if shouldBeEqual {
						t.Fatalf("Mismatch in field %q: got %v, expected %v", field, got, expected)
					} else {
						t.Fatalf("Field %q should have changed, but it did not: got %v", field, got)
					}
				},
				testutils.WithFieldNotEqual[model.User]("UpdatedAt"),
				testutils.WithFieldNotEqual[model.User]("Email"),
				testutils.WithFieldNotEqual[model.User]("Password"),
			)
		})

		t.Run("no fields", func(t *testing.T) {
			t.Parallel()

			createdUser, err := repo.Create(testUsers[1].email, testUsers[1].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			_, err = repo.Update(createdUser.ID)
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for missing fields")
		})

		t.Run("update email", func(t *testing.T) {
			t.Parallel()

			createdUser, err := repo.Create(testUsers[3].email, testUsers[3].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			updatedUser, err := repo.Update(createdUser.ID,
				model.WithEmail("updated-email-3"))
			assert.NoErrorf(t, err, "Update User failed: %v", err)

			testutils.AssertEqualItems(createdUser, updatedUser,
				func(field string, shouldBeEqual bool, got any, expected any) {
					t.Helper()
					if shouldBeEqual {
						t.Fatalf("Mismatch in field %q: got %v, expected %v", field, got, expected)
					} else {
						t.Fatalf("Field %q should have changed, but it did not: got %v", field, got)
					}
				},
				testutils.WithFieldNotEqual[model.User]("UpdatedAt"),
				testutils.WithFieldNotEqual[model.User]("Email"),
			)
		})

		t.Run("update password", func(t *testing.T) {
			t.Parallel()

			createdUser, err := repo.Create(testUsers[2].email, testUsers[2].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			updatedUser, err := repo.Update(createdUser.ID,
				model.WithPassword("updated-password-2"))
			assert.NoErrorf(t, err, "Update User failed: %v", err)

			testutils.AssertEqualItems(createdUser, updatedUser,
				func(field string, shouldBeEqual bool, got any, expected any) {
					t.Helper()
					if shouldBeEqual {
						t.Fatalf("Mismatch in field %q: got %v, expected %v", field, got, expected)
					} else {
						t.Fatalf("Field %q should have changed, but it did not: got %v", field, got)
					}
				},
				testutils.WithFieldNotEqual[model.User]("UpdatedAt"),
				testutils.WithFieldNotEqual[model.User]("Password"),
			)
		})
	})

	t.Run("delete user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := repositories.NewSQLUserRepository(db)

		t.Run("by id", func(t *testing.T) {
			t.Parallel()

			user, err := repo.Create(testUsers[0].email, testUsers[0].password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)

			retrieved, err := repo.GetByID(user.ID)
			assertGetUser(t, user, retrieved, err)

			err = repo.Delete(user.ID)
			assert.NoErrorf(t, err, "Delete User failed: %v", err)

			_, err = repo.GetByID(user.ID)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error not found")
		})

		t.Run("by id - none", func(t *testing.T) {
			t.Parallel()

			err := repo.Delete(model.ID(99999))
			assert.NoErrorf(t, err, "Delete User failed: %v", err)
		})
	})
}

func generateTestUsers() []testUser {
	const size = 10
	testUsers := make([]testUser, size)

	for i := range testUsers {
		testUsers[i] = testUser{
			id:       uint(i),
			email:    fmt.Sprintf("user-%d@dwimc.awesome", i),
			password: fmt.Sprintf("secret-password-%d", i),
		}
	}

	return testUsers
}

func assertCreateUser(t *testing.T, expected *testUser, actual *model.User, err error) {
	assert.NoErrorf(t, err, "Create User failed: %v", err)

	assert.NotEqualf(t, 0, actual.ID, "Invalid ID 0")
	assert.Equalf(t, expected.email, actual.Email, "Email mismatch (user ID: %d)", expected.id)
	assert.Equalf(t, expected.password, actual.Password, "Password mismatch (user ID: %d)", expected.id)
}

func assertGetUser(t *testing.T, actual *model.User, expected *model.User, err error) {
	assert.NoErrorf(t, err, "Get User failed: %v", err)
	assert.Equalf(t, expected, actual, "Got %+v, expected: %+v (user ID: %d)", actual, expected, expected.ID)
}
