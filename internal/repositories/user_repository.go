package repositories

import (
	"database/sql"
	"dwimc/internal/model"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type userUpdateField struct{}

func (userUpdateField) WithPassword(password string) UpdateField {
	return WithField("password", password)
}

func (userUpdateField) WithToken(token string) UpdateField {
	return WithField("token", token)
}

var UserUpdate userUpdateField


type UserRepository interface {
	GetBy(id model.ID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(email string, password string, token string) (*model.User, error)
	Update(id model.ID, fields ...UpdateField) (*model.User, error)
	Delete(id model.ID) error
}

type SQLUserRepository struct {
	db *sqlx.DB
}

func NewSQLUserRepository(db *sqlx.DB) UserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetBy(id model.ID) (*model.User, error) {
	query := `
		SELECT id, created_at, updated_at,
				email, password, token
		FROM users
		WHERE id = ?
	`

	return getUserBy(r.db, query, id)
}

func (r *SQLUserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, created_at, updated_at,
				email, password, token
		FROM users
		WHERE email = ?
	`

	return getUserBy(r.db, query, email)
}

func (r *SQLUserRepository) Create(email string, password string, token string) (*model.User, error) {
	query := `
		INSERT INTO users (email, password, token)
			VALUES ($1, $2, $3)
			RETURNING *
	`
	user := model.User{}

	err := r.db.Get(&user, query, email, password, token)
	// TODO - Handle constraints (email unique) errors - User Already Exist
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *SQLUserRepository) Update(id model.ID, fields ...UpdateField) (*model.User, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("Update error: missing fields")
	}

	query := "UPDATE users SET "
	updates := map[string]interface{}{}
	setClauses := []string{}
	args := []interface{}{}

	for _, field := range fields {
		field(&updates)
	}

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	query += strings.Join(setClauses, ", ") + " WHERE id = ? RETURNING *"
	args = append(args, id)

	updatedUser := model.User{}

	err := r.db.Get(&updatedUser, query, args...)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *SQLUserRepository) Delete(id model.ID) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		// TODO - handle errors?
		return err
	}

	return nil
}

func getUserBy[T model.ID | string](db *sqlx.DB, query string, field T) (*model.User, error) {
	user := model.User{}
	err := db.Get(&user, query, field)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
