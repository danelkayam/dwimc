package repositories

import (
	"dwimc/internal/model"

	"github.com/jmoiron/sqlx"
)

type LocationRepository interface {
	GetLast(deviceID model.ID) (*model.Location, error)
	GetAllBy(deviceID model.ID) ([]model.Location, error)
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

func (r *SQLLocationRepository) GetLast(deviceID model.ID) (*model.Location, error) {
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
		return nil, handleSQLError("failed getting location", err)
	}

	return &location, nil
}

func (r *SQLLocationRepository) GetAllBy(deviceID model.ID) ([]model.Location, error) {
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
		return nil, handleSQLError("failed getting locations by device", err)
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
		return nil, handleSQLError("failed creating location", err)
	}

	return &location, nil
}

func (r *SQLLocationRepository) Delete(id model.ID) error {
	query := `DELETE FROM locations WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return handleSQLError("failed deleting location", err)
	}

	return nil
}

func (r *SQLLocationRepository) DeleteAllBy(deviceID model.ID) (int64, error) {
	query := `DELETE FROM locations WHERE device_id = ?`

	res, err := r.db.Exec(query, deviceID)
	if err != nil {
		return 0, handleSQLError("failed deleting locations by device", err)
	}

	total, err := res.RowsAffected()
	if err != nil {
		return 0, handleSQLError("failed deleting locations by device", err)
	}

	return total, nil
}
