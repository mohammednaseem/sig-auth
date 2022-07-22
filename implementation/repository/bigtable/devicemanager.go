package btable

import (
	bigtable "cloud.google.com/go/bigtable"
	"github.com/device-auth/model"
)

type deviceRepository struct {
	Conn *bigtable.Client
	Table *bigtable.Table
}

func NewDeviceRepository(Conn *bigtable.Client,Table *bigtable.Table) model.IDeviceRepository {
	return &deviceRepository{Conn: Conn,Table:Table}
}
