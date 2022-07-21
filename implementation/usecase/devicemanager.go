package usecase

import (
	"time"

	"github.com/device-auth/model"
)

type dDeviceUsecase struct {
	deviceRepo     model.IDeviceRepository
	contextTimeout time.Duration
}

func NewDeviceUsecase(d model.IDeviceRepository, timeout time.Duration) model.IDeviceUsecase {
	return &dDeviceUsecase{
		deviceRepo:     d,
		contextTimeout: timeout,
	}
}
