package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	Service "github.com/device-auth/implementation/service"
	"github.com/device-auth/model"
	jwt "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

func CreateDevicePublish(topicId string, dev model.DeviceCreate) error {

	PubStruct := model.PublishDeviceCreate{Operation: "POST", Entity: "Device", Data: dev, Path: "device/" + dev.Parent}

	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Service.Publish(dev.Project, topicId, msg)

	return err
}
func GetDeviceData(deviceId string, tokenString string) (model.DeviceCreate, error) {
	var signingMethod string
	var dev model.DeviceCreate
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			signingMethod = "ES256"
		} else if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			signingMethod = "RS256"
		} else if signingMethod == "" {
			return false, errors.New("unknown signing method")
		}
		return true, nil
	})
	dev.Project = fmt.Sprintf("%v", claims["aud"])
	PubStruct := model.PublishDeviceCreate{Operation: "POST", Entity: "Device", Data: dev, Path: "device/" + dev.Parent}

	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	err = Service.Publish(dev.Project, topicId, msg)

	return err
}
func (d *dDeviceUsecase) CheckCredentials(ctx context.Context, input model.DeviceAndToken) (bool, error) {
	var Certs []string
	var err error
	if input.Bootstrap != "" {
		Certs, err = d.GetCertificateFromDb(ctx, input.DeviceId, input.Token)
		if err != nil {
			log.Error().Err(err).Msg("")
			return false, err
		}
	} else {
		Certs = append(Certs, input.Bootstrap)
	}

	boolVal, err := d.IsCertificateKeyMapped(ctx, Certs, input.Token)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	err = GetDeviceData(input.DeviceId, input.Token)
	err = CreateDevicePublish(d.topicId, dev)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	return boolVal, err
}
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

func (d *dDeviceUsecase) GetCertificateFromDb(ctx context.Context, deviceId string, token string) ([]string, error) {
	c, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	mDevice, err := d.GetAllPublicKeysForDevice(c, deviceId)
	var Certs []string
	if err != nil {
		log.Error().Err(err).Msg("")
		return Certs, err
	}

	for _, element := range mDevice.Credentials {
		if len(strings.TrimSpace(element.PublicKey.Key)) != 0 {
			Certs = append(Certs, element.PublicKey.Key)
		}
	}

	return Certs, err
}
func (d *dDeviceUsecase) IsCertificateKeyMapped(ctx context.Context, certificate []string, token string) (bool, error) {

	isValidDevice, err := IdentifyAndVerifyJWT(token, certificate)

	return isValidDevice, err
}
