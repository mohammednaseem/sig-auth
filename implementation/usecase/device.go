package usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/device-auth/model"
)

func (d *dDeviceUsecase) GetAllPublicKeysForDevice(ctx context.Context, deviceId string) (model.Device, error) {
	c, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	mDevice, err := d.deviceRepo.GetAllPublicKeysForDevice(c, deviceId)
	if err != nil {
		log.Error().Err(err).Msg("")
		return model.Device{}, err
	}
	return mDevice, err
}

func (d *dDeviceUsecase) IsValidCertificate(ctx context.Context, deviceId string, token string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	mDevice, err := d.GetAllPublicKeysForDevice(c, deviceId)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	isValidDevice, err := IdentifyAndVerifyJWT(token, mDevice)

	return isValidDevice, err
}
