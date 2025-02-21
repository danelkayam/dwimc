package repositories

import (
	"database/sql"
	"dwimc/internal/model"
	"errors"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetBy(id model.ID) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(email string, password string, token string) (*model.User, error)
	Update(id model.ID, password string, token string) (*model.User, error)
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
	newUser := model.User{}

	err := r.db.Get(&newUser, query, email, password, token)
	// TODO - Handle constraints (email unique) errors - User Already Exist
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *SQLUserRepository) Update(id model.ID, password string, token string) (*model.User, error) {
	// TODO: handle updating only password or token, if both empty, throw error
	query := `
		UPDATE users
			SET password = ?,
				token = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
			RETURNING *
	`
	updatedUser := model.User{}

	err := r.db.Get(&updatedUser, query, token, password, id)
	if err != nil {
		// TODO - handle errors
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
			// TODO - handle this better?
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
