package usecase

import (
	"time"

	"github.com/device-auth/model"
)

type dDeviceUsecase struct {
	deviceRepo     model.IDeviceRepository
	contextTimeout time.Duration
	topicId        string
}

func NewDeviceUsecase(d model.IDeviceRepository, timeout time.Duration, topicId string) model.IDeviceUsecase {
	return &dDeviceUsecase{
		deviceRepo:     d,
		contextTimeout: timeout,
		topicId:        topicId,
	}
}
