package mongo

import (
	"context"

	"github.com/device-auth/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type deviceRepository struct {
	client      *mongo.Client
	rcollection string
	dcollection string
	database    string
	ctx         context.Context
}

func NewDeviceRepository(ctx context.Context, conn *mongo.Client, dcollection string, rcollection string, database string) model.IDeviceRepository {
	return &deviceRepository{
		client:      conn,
		rcollection: rcollection,
		dcollection: dcollection,
		database:    database,
		ctx:         ctx}
}
