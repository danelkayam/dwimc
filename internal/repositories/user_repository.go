package repositories

import (
	"dwimc/internal/model"
	"dwimc/internal/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByID(id model.ID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(email string, password string) (*model.User, error)
	Update(id model.ID, fields ...model.UpdateField) (*model.User, error)
	Delete(id model.ID) error
}

type SQLUserRepository struct {
	db *sqlx.DB
}

func NewSQLUserRepository(db *sqlx.DB) UserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetByID(id model.ID) (*model.User, error) {
	query := `
		SELECT id, created_at, updated_at, email, password
		FROM users
		WHERE id = ?
	`

	return getUserBy(r.db, query, id)
}

func (r *SQLUserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, created_at, updated_at, email, password
		FROM users
		WHERE email = ?
	`

	return getUserBy(r.db, query, email)
}

func (r *SQLUserRepository) Create(email string, password string) (*model.User, error) {
	query := `
		INSERT INTO users (email, password)
			VALUES ($1, $2)
			RETURNING *
	`
	var user model.User

	err := r.db.Get(&user, query, email, password)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return &user, nil
}

func (r *SQLUserRepository) Update(id model.ID, fields ...model.UpdateField) (*model.User, error) {
	if len(fields) == 0 {
		return nil, utils.AsError(model.ErrInvalidArgs, "missing fields")
	}

	query := "UPDATE users SET "
	updates := map[string]any{}
	setClauses := []string{}
	args := []any{}

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

	var user model.User

	err := r.db.Get(&user, query, args...)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return &user, nil
}

func (r *SQLUserRepository) Delete(id model.ID) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return handleSQLError(err)
	}

	return nil
}

func getUserBy[T model.ID | string](db *sqlx.DB, query string, field T) (*model.User, error) {
	var user model.User

	err := db.Get(&user, query, field)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return &user, nil
}
