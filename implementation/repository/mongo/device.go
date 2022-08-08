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
		{Key: "name", Value: bson.D{{Key: "$eq", Value: deviceId}}},
	}

	// Returns result of deletion and error
	err = queryOne(ctx, client, db, collection, filter).Decode(&mDevice)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	// print the count of affected documents
	if mDevice.Id == "" {
		log.Error().Msg("No Data Found in Db")
	}
	log.Info().Msg("Got Details For Device" + mDevice.Id)
	return
}

func (d *deviceRepository) GetAllPublicKeysForDevice(_ context.Context, deviceId string) (model.Device, error) {

	mDevices, err := getDeviceDetails(d.ctx, d.client, d.database, d.dcollection, deviceId)

	return mDevices, err
}
func (d *deviceRepository) CheckCaSign(_ context.Context, registry string, region string, project string, bootstrap string) (bool, error) {
	Ping(d.ctx, d.client)
	var filter interface{} = bson.D{
		{Key: "id", Value: bson.D{{Key: "$eq", Value: registry}}},
		{Key: "region", Value: bson.D{{Key: "$eq", Value: region}}},
		{Key: "project", Value: bson.D{{Key: "$eq", Value: project}}},
	}
	var mDevice model.Registry
	var err error
	// Returns result of deletion and error
	err = queryOne(d.ctx, d.client, d.database, d.rcollection, filter).Decode(&mDevice)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	// print the count of affected documents
	if mDevice.Id == "" {
		log.Error().Msg("No Data Found in Db")
		return false, err
	}
	if len(mDevice.Credentials) > 0 {

		for _, ca := range mDevice.Credentials {
			err = verifyCert([]byte(bootstrap), []byte(ca.PublicKeyCertificate.Certificate))
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Error().Msg("Certificate Verification Failed")
			return false, err
		}
	}
	log.Info().Msg("Got Details For Device" + mDevice.Id)

	return true, err
}
