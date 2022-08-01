package mongo

import (
	"context"

	"github.com/device-auth/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getDeviceDetails(ctx context.Context, client *mongo.Client, db string, collection string, deviceId string) (mDevice model.Device, err error) {
	Ping(ctx, client)
	var filter interface{} = bson.D{
		{Key: "id", Value: bson.D{{Key: "$eq", Value: deviceId}}},
	}

	// Returns result of deletion and error
	var queryResult model.Device
	err = queryOne(ctx, client, db, collection, filter).Decode(&queryResult)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	// print the count of affected documents
	if queryResult.Id == "" {
		log.Error().Msg("No Data Found in Db")
	}
	log.Info().Msg("Got Details For Device" + queryResult.Id)
	return
}

func (d *deviceRepository) GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (model.Device, error) {

	mDevices, err := getDeviceDetails(d.ctx, d.client, d.database, d.collection, deviceId)

	return mDevices, err
}
