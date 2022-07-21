package psql

import (
	"database/sql"

	"github.com/device-auth/model"
)

type deviceRepository struct {
	Conn *sql.DB
}

func NewDeviceRepository(Conn *sql.DB) model.IDeviceRepository {
	return &deviceRepository{Conn: Conn}
}
