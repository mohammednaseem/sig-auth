package http

import (
	"github.com/device-auth/model"
	"github.com/labstack/echo"
)

type ResponseError struct {
	Message interface{}
}

type deviceHandler struct {
	dUsecase model.IDeviceUsecase
}

func NewDeviceHandler(e *echo.Echo, dUsecase model.IDeviceUsecase) {
	DeviceHandler := &deviceHandler{
		dUsecase: dUsecase,
	}
	e.GET("/device", DeviceHandler.getDeviceDetails)
	e.POST("/device/authenticate", DeviceHandler.authN)
}
