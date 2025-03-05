package repositories

import (
	"database/sql"
	"dwimc/internal/model"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type deviceUpdateField struct{}

func (deviceUpdateField) WithSerial(serial string) UpdateField {
	return WithField("serial", serial)
}

func (deviceUpdateField) WithName(name string) UpdateField {
	return WithField("name", name)
}

func (deviceUpdateField) WithToken(token string) UpdateField {
	return WithField("token", token)
}

var DeviceUpdate deviceUpdateField

type DeviceRepository interface {
	Get(id model.ID) (*model.Device, error)
	GetBySerial(userID model.ID, serial string) (*model.Device, error)
	GetAllByUserID(userID model.ID) ([]model.Device, error)
	Create(userId model.ID, serial string, name string, token string) (*model.Device, error)
	Update(id model.ID, fields ...UpdateField) (*model.Device, error)
	Delete(id model.ID) error
	DeleteAllByUserID(userID model.ID) (int64, error)
}

type SQLDeviceRepository struct {
	db *sqlx.DB
}

func NewSQLDeviceRepository(db *sqlx.DB) DeviceRepository {
	return &SQLDeviceRepository{db: db}
}

func (r *SQLDeviceRepository) Get(id model.ID) (*model.Device, error) {
	query := `
		SELECT id, created_at, updated_at,
			user_id, serial, name, token
		FROM devices
		WHERE id = ?
	`

	var device model.Device
	err := r.db.Get(&device, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &device, nil
}

func (r *SQLDeviceRepository) GetBySerial(userID model.ID, serial string) (*model.Device, error) {
	query := `
		SELECT id, created_at, updated_at,
			user_id, serial, name, token
		FROM devices
		WHERE user_id = ? AND serial = ?
	`

	var device model.Device
	err := r.db.Get(&device, query, userID, serial)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &device, nil
}

func (r *SQLDeviceRepository) GetAllByUserID(userID model.ID) ([]model.Device, error) {
	query := `
		SELECT id, created_at, updated_at,
			user_id, serial, name, token
		FROM devices
		WHERE user_id = ?
		ORDER BY created_at ASC
	`

	devices := []model.Device{}
	err := r.db.Select(&devices, query, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return devices, nil
		}

		return nil, err
	}

	return devices, nil
}

func (r *SQLDeviceRepository) Create(userID model.ID, serial string, name string, token string) (*model.Device, error) {
	query := `
	INSERT INTO devices (user_id, serial, name, token)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	var device model.Device
	err := r.db.Get(&device, query, userID, serial, name, token)
	
	// TODO - handle constraints
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *SQLDeviceRepository) Update(id model.ID, fields ...UpdateField) (*model.Device, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("Update error: missing fields")
	}

	query := "UPDATE devices SET "
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

	var device model.Device
	err := r.db.Get(&device, query, args...)
	
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (r *SQLDeviceRepository) Delete(id model.ID) error {
	query := `DELETE FROM devices WHERE id = ?`

	_, err := r.db.Exec(query, id)
	
	if err != nil {
		// TODO - handle errors?
		return err
	}

	return nil
}

func (r *SQLDeviceRepository) DeleteAllByUserID(userID model.ID) (int64, error) {
	query := `DELETE FROM devices WHERE user_id = ?`

	res, err := r.db.Exec(query, userID)
	
	if err != nil {
		// TODO - handle errors?
		return 0, err
	}

	total, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return total, nil
}
