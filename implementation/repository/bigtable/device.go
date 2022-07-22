package btable

import (
	"context"
	"log"
	"fmt"
	"strings"
	bigtable "cloud.google.com/go/bigtable"
	"github.com/device-auth/model"
	Log "github.com/rs/zerolog/log"
)
func printRow(row bigtable.Row) {
	fmt.Printf("Reading data for %s:\n", row.Key())
	for columnFamily, cols := range row {
			fmt.Printf("Column Family %s\n", columnFamily)
			for _, col := range cols {
					qualifier := col.Column[strings.IndexByte(col.Column, ':')+1:]
					fmt.Printf("\t%s: %s @%d\n", qualifier, col.Value, col.Timestamp)
			}
	}
}
func getDeviceDetails(ctx context.Context, db *bigtable.Client,table *bigtable.Table, query string, deviceId string) (mDevice model.Device, err error) {
	Log.Debug().Str("query", query).Msg("")
	row, err := table.ReadRow(ctx, deviceId)
	if err != nil {
			log.Fatalf("Could not read row with key %s: %v", deviceId, err)
	}
	printRow(row)
	return
}

func (d *deviceRepository) GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (model.Device, error) {
	query := `SELECT * FROM public.device  WHERE deviceid=$1;`
	mDevices, err := getDeviceDetails(ctx, d.Conn,d.Table, query, deviceId)
	if err != nil {
		log.Fatalf("Error Fetching From Big Table", err)
	}	

	
	return mDevices, nil
}
