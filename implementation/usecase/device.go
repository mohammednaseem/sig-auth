package usecase

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

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
	nBig, err := rand.Int(rand.Reader, big.NewInt(999999999999999999))
	if err != nil {
		log.Error().Msg("Random Generator Failed")
		return model.DeviceCreate{}, err
	}
	randNum := nBig.Int64()
	dev.NumId = fmt.Sprintf("%d", randNum)
	dev.Name = deviceId
	devInfo := strings.Split(deviceId, "/")
	if len(devInfo) != 8 {
		err := errors.New("mqtt clientId unknown format")
		log.Error().Err(err).Msg("")
		return model.DeviceCreate{}, err
	}
	dev.CreatedOn = time.Now().String()

	dev.Project = devInfo[1]
	dev.Region = devInfo[3]
	dev.Registry = devInfo[5]
	dev.Id = devInfo[7]
	dev.Parent = fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices", dev.Project, dev.Region, dev.Registry)
	dev.Blocked = false
	dev.Metadata = map[string]string{}
	dev.LogLevel = "INFO"
	s, _ := json.MarshalIndent(dev, "", "\t")
	fmt.Print(string(s))

	return dev, err
}
func (d *dDeviceUsecase) CheckCredentials(ctx context.Context, input model.DeviceAndToken) (bool, error) {
	var Certs []string
	var err error
	var zeroTouch bool = true
	block, _ := pem.Decode([]byte(input.Bootstrap))
	if block == nil {
		log.Info().Msg("Bootstrap Certificate Invalid")
		zeroTouch = false
	}
	if !zeroTouch {
		Certs, err = d.GetCertificateFromDb(ctx, input.DeviceId, input.Token)
		if err != nil {
			log.Error().Err(err).Msg("")
			return false, err
		}
	} else {
		Certs = append(Certs, input.Bootstrap)
	}
	if len(Certs) == 0 {
		log.Error().Msg("Certificates Not Present For Check")
		return false, errors.New("certificates not present for check")
	}
	boolVal, err, algorithm := d.IsCertificateKeyMapped(ctx, Certs, input.Token)

	if err != nil {
		log.Error().Err(err).Msg("")
		return false, err
	}
	log.Info().Msg("Token Verified")
	if zeroTouch {
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
