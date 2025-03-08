package repositories

import (
	"dwimc/internal/model"

	"github.com/jmoiron/sqlx"
)

type LocationRepository interface {
	GetLastByDeviceID(deviceID model.ID) (*model.Location, error)
	GetAllByDeviceID(deviceID model.ID) ([]model.Location, error)
	Create(deviceID model.ID, latitude float64, longitude float64) (*model.Location, error)
	Delete(id model.ID) error
	DeleteAllBy(deviceID model.ID) (int64, error)
}

type SQLLocationRepository struct {
	db *sqlx.DB
}

func NewSQLLocationRepository(db *sqlx.DB) LocationRepository {
	return &SQLLocationRepository{db: db}
}

func (r *SQLLocationRepository) GetLastByDeviceID(deviceID model.ID) (*model.Location, error) {
	query := `
		SELECT id, created_at, updated_at,
			device_id, latitude, longitude
		FROM locations
		WHERE device_id = ?
		ORDER BY created_at DESC
		LIMIT 1;
	`

	var location model.Location

	err := r.db.Get(&location, query, deviceID)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return &location, nil
}

func (r *SQLLocationRepository) GetAllByDeviceID(deviceID model.ID) ([]model.Location, error) {
	query := `
		SELECT id, created_at, updated_at,
			device_id, latitude, longitude
		FROM locations
		WHERE device_id = ?
		ORDER BY created_at DESC
	`

	locations := []model.Location{}

	err := r.db.Select(&locations, query, deviceID)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return locations, nil
}

func (r *SQLLocationRepository) Create(deviceID model.ID, latitude float64, longitude float64) (*model.Location, error) {
	query := `
		INSERT INTO locations (device_id, latitude, longitude)
			VALUES ($1, $2, $3)
			RETURNING *
	`

	var location model.Location

	err := r.db.Get(&location, query, deviceID, latitude, longitude)
	if err != nil {
		return nil, handleSQLError(err)
	}

	return &location, nil
}

func (r *SQLLocationRepository) Delete(id model.ID) error {
	query := `DELETE FROM locations WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return handleSQLError(err)
	}

	return nil
}

func (r *SQLLocationRepository) DeleteAllBy(deviceID model.ID) (int64, error) {
	query := `DELETE FROM locations WHERE device_id = ?`

	res, err := r.db.Exec(query, deviceID)
	if err != nil {
		return 0, handleSQLError(err)
	}

	total, err := res.RowsAffected()
	if err != nil {
		return 0, handleSQLError(err)
	}

	return total, nil
}
