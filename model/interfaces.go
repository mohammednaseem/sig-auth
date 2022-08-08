package model

import (
	"context"
)

//jwt usecase
type IDeviceUsecase interface {
	CheckCredentials(ctx context.Context, deviceId DeviceAndToken) (bool, error)
	GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (Device, error)
	GetCertificateFromDb(ctx context.Context, deviceId string) ([]string, error)
	IsCertificateKeyMapped(ctx context.Context, certificate []string, token string) (bool, string, error)
}

//jwt repo
type IDeviceRepository interface {
	GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (Device, error)
	CheckCaSign(ctx context.Context, registry string, region string, project string, bootstrap string) (bool, error)
}
