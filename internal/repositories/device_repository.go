package repositories

import (
	"dwimc/internal/model"

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
	GetAllByUserID(userID model.ID) ([]model.Device, error)
	Create(userId model.ID, serial string, name string, token string) (*model.Device, error)
	Update(id model.ID, fields ...UpdateField) (*model.Device, error)
	Delete(id model.ID) error
	DeleteAllByUserID(userID model.ID) (int, error)
}

type SQLDeviceRepository struct {
	db *sqlx.DB
}

func NewSQLDeviceRepository(db *sqlx.DB) DeviceRepository {
	return &SQLDeviceRepository{db: db}
}

func (r *SQLDeviceRepository) Get(id model.ID) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *SQLDeviceRepository) GetAllByUserID(userID model.ID) ([]model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *SQLDeviceRepository) Create(userId model.ID, serial string, name string, token string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *SQLDeviceRepository) Update(id model.ID, fields ...UpdateField) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (r *SQLDeviceRepository) Delete(id model.ID) error {
	// TODO - implement this
	return nil
}

func (r *SQLDeviceRepository) DeleteAllByUserID(id model.ID) (int, error) {
	// TODO - implement this
	return 0, nil
}
