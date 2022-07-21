package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/device-auth/helper"
	"github.com/device-auth/model"
	"github.com/rs/zerolog/log"
)

func getDeviceDetails(ctx context.Context, db *sql.DB, query string, deviceId string) (mDevice model.Device, err error) {
	log.Debug().Str("query", query).Msg("")
	row := db.QueryRow(query, deviceId)
	if err != nil {
		err = helper.CheckDatabaseError(err)
		fmt.Println(err)
		return model.Device{}, err
	}
	var Cerificate1, Cerificate2, Cerificate3 sql.NullString
	switch err := row.Scan(&mDevice.DeviceId, &mDevice.Name, &mDevice.Password, &Cerificate1, &Cerificate2, &Cerificate3, &mDevice.Project, &mDevice.Region, &mDevice.Created_On); err {
	case sql.ErrNoRows:
		fmt.Print("There is no retrieved rows, dummy!")
	case nil:
		fmt.Print(mDevice.DeviceId, "\n")
	default:
		panic(err)
	}
	if Cerificate1.Valid {
		mDevice.Cerificate1 = Cerificate1.String
	}
	if Cerificate2.Valid {
		mDevice.Cerificate2 = Cerificate2.String
	}
	if Cerificate3.Valid {
		mDevice.Cerificate3 = Cerificate3.String
	}
	return
}

func (d *deviceRepository) GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (model.Device, error) {
	query := `SELECT * FROM public.device  WHERE deviceid=$1;`
	mDevices, err := getDeviceDetails(ctx, d.Conn, query, deviceId)

	switch {
	case err == sql.ErrNoRows:
		return mDevices, nil
	case err != nil:
		err = helper.CheckDatabaseError(err)
		return mDevices, err
	}
	return mDevices, nil
}
