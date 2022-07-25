package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/device-auth/helper"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func (d *deviceHandler) getDeviceDetails(c echo.Context) error {
	ctx := c.Request().Context()
	deviceid := c.QueryParam("deviceid")
	if len(strings.TrimSpace(deviceid)) == 0 {
		r := ResponseError{Message: "Missing query param DeviceId"}
		return c.JSON(http.StatusBadRequest, r)
	}

	_, err := strconv.Atoi(deviceid)
	if err != nil {
		r := ResponseError{Message: "Query param DeviceId should be an integer"}
		return c.JSON(http.StatusBadRequest, r)
	}

	device, err := d.dUsecase.GetAllPublicKeysForDevice(ctx, deviceid)
	if err != nil {
		log.Error().Err(err).Str("Method", "GetEnvironment").Int("Environment", 1).Msg("")
		return c.JSON(helper.GetStatusCode(err), ResponseError{Message: err.Error()})
	}
	if len(strings.TrimSpace(device.DeviceId)) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	} else {
		return c.JSON(http.StatusOK, device)
	}
}

type DeviceAndToken struct {
	DeviceId string `json:"deviceid" validate:"required"`
	Token    string `json:"token" validate:"required,min=10"`
}
type ResponseSuccess struct {
	Result       string `json:"result" xml:"result"`
	Is_superuser bool   `json:"is_superuser" xml:"is_superuser"`
}

func (d *deviceHandler) authN(c echo.Context) error {
	ctx := c.Request().Context()

	dt := new(DeviceAndToken)
	if err := c.Bind(dt); err != nil {
		log.Error().Err(err).Msg("")
		r := ResponseError{Message: "Data not good"}
		return c.JSON(http.StatusBadRequest, r)
	}
	boolVal, err := d.dUsecase.IsValidCertificate(ctx, dt.DeviceId, dt.Token)
	if !boolVal {
		return c.JSON(http.StatusForbidden, err)
	} else {
		r := ResponseSuccess{Result: "allow", Is_superuser: false}
		return c.JSON(http.StatusOK, r)
	}

}
