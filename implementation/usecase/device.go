package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	Service "github.com/device-auth/implementation/service"
	"github.com/device-auth/model"
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
func GetDeviceData(deviceId string, tokenString string, algorithm string, topicId string) (model.DeviceCreate, error) {
	var dev model.DeviceCreate
	//claims := jwt.MapClaims{}
	//_, err := jwt.ParseWithClaims(tokenString, claims, nil)
	//dev.Project = fmt.Sprintf("%v", claims["aud"])
	dev.Name = deviceId
	fmt.Sscanf(deviceId, "projects/%v/locations/%v/registries/%v/devices/%v", &dev.Project, &dev.Region, &dev.Registry, &dev.Id)
	dev.Blocked = false
	dev.Metadata = map[string]string{}
	dev.LogLevel = "INFO"
	PubStruct := model.PublishDeviceCreate{Operation: "POST", Entity: "Device", Data: dev, Path: "device/" + dev.Parent}
	s, _ := json.MarshalIndent(dev, "", "\t")
	fmt.Print(string(s))
	msg, err := json.Marshal(PubStruct)
	if err != nil {
		log.Error().Err(err).Msg("")
		return model.DeviceCreate{}, err
	}
	err = Service.Publish(dev.Project, topicId, msg)

	return dev, err
}
func (d *dDeviceUsecase) CheckCredentials(ctx context.Context, input model.DeviceAndToken) (bool, error) {
	var Certs []string
	var err error
	if input.Bootstrap == "" {
		Certs, err = d.GetCertificateFromDb(ctx, input.DeviceId, input.Token)
		if err != nil {
			log.Error().Err(err).Msg("")
			return false, err
		}
	} else {
		Certs = append(Certs, input.Bootstrap)
	}

	boolVal, err, algorithm := d.IsCertificateKeyMapped(ctx, Certs, input.Token)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	dev, err := GetDeviceData(input.DeviceId, input.Token, algorithm, d.topicId)
	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
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
func (d *dDeviceUsecase) IsCertificateKeyMapped(ctx context.Context, certificate []string, token string) (bool, error, string) {

	isValidDevice, err, algorithm := IdentifyAndVerifyJWT(token, certificate)

	return isValidDevice, err, algorithm
}
