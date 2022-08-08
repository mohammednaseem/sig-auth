package usecase

import (
	"time"

	"github.com/device-auth/model"
)

type dDeviceUsecase struct {
	deviceRepo     model.IDeviceRepository
	contextTimeout time.Duration
	topicId        string
	projectId      string
}

func NewDeviceUsecase(d model.IDeviceRepository, timeout time.Duration, topicId string, projectId string) model.IDeviceUsecase {
	return &dDeviceUsecase{
		deviceRepo:     d,
		contextTimeout: timeout,
		topicId:        topicId,
		projectId:      projectId,
	}
}
