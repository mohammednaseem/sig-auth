package model

import (
	"context"
)

//jwt usecase
type IDeviceUsecase interface {
	GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (Device, error)
	IsValidCertificate(ctx context.Context, deviceId string, token string) (bool, error)
}

//jwt repo
type IDeviceRepository interface {
	GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (Device, error)
}
