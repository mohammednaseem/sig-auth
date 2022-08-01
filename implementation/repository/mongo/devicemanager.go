package mongo

import (
	"context"

	"github.com/device-auth/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type deviceRepository struct {
	client     *mongo.Client
	collection string
	database   string
	ctx        context.Context
}

func NewDeviceRepository(ctx context.Context, conn *mongo.Client, collection string, database string) model.IDeviceRepository {
	return &deviceRepository{
		client:     conn,
		collection: collection,
		database:   database,
		ctx:        ctx}
}
