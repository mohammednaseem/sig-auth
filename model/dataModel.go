package model

import "google.golang.org/api/cloudiot/v1"

type Device struct {
	Project     string                       `json:"project" validate:"required"`
	Parent      string                       `json:"parent" validate:"required"`
	NumId       string                       `json:"numId" validate:""`
	Region      string                       `json:"region" validate:"required"`
	Registry    string                       `json:"registry" validate:"required"`
	Id          string                       `json:"id" validate:"required"`
	Name        string                       `json:"name" validate:"required"`
	Credentials []*cloudiot.DeviceCredential `json:"credentials" validate:"required"`
	LogLevel    string                       `json:"logLevel"  validate:""`
	Blocked     bool                         `json:"blocked"  validate:""`
	Metadata    map[string]string            `json:"metadata"  validate:""`
	CreatedOn   string                       `json:"createdOn"  validate:""`
}
type DeviceAndToken struct {
	DeviceId  string `json:"deviceid" validate:"required"`
	Token     string `json:"token" validate:"required,min=10"`
	Bootstrap string `json:"bootstrap" validate:"required,min=10"`
}
type DeviceCreate struct {
	Project       string                       `json:"project" validate:"required"`
	Parent        string                       `json:"parent" validate:"required"`
	NumId         string                       `json:"numId" validate:""`
	Region        string                       `json:"region" validate:"required"`
	Registry      string                       `json:"registry" validate:"required"`
	Id            string                       `json:"id" validate:"required"`
	Name          string                       `json:"name" validate:"required"`
	Credentials   []*cloudiot.DeviceCredential `json:"credentials" validate:"required"`
	LogLevel      string                       `json:"logLevel"  validate:""`
	Blocked       bool                         `json:"blocked"  validate:""`
	Metadata      map[string]string            `json:"metadata"  validate:""`
	CreatedOn     string                       `json:"createdOn"  validate:""`
	Decomissioned bool                         `json:"decomissioned"  validate:""`
}
type PublishDeviceCreate struct {
	Operation string       `json:"operation" validate:"required"`
	Entity    string       `json:"entity" validate:"required"`
	Path      string       `json:"path" validate:"required"`
	Data      DeviceCreate `json:"data" validate:"required"`
}
type Registry struct {
	Parent                   string                              `json:"parent" validate:"required"`
	Project                  string                              `json:"project" validate:"required"`
	Region                   string                              `json:"region" validate:"required"`
	Id                       string                              `json:"id" validate:"required"`
	Name                     string                              `json:"name" validate:"required"`
	EventNotificationConfigs []*cloudiot.EventNotificationConfig `json:"eventNotificationConfigs" validate:"required"`
	StateNotificationConfig  *cloudiot.StateNotificationConfig   `json:"stateNotificationConfig"  validate:""`
	MqttConfig               cloudiot.MqttConfig                 `json:"mqttConfig"  validate:""`
	HttpConfig               cloudiot.HttpConfig                 `json:"httpConfig"  validate:""`
	LogLevel                 string                              `json:"logLevel"  validate:""`
	CreatedOn                string                              `json:"createdOn"  validate:""`
	Credentials              []*cloudiot.RegistryCredential      `json:"credentials"  validate:""`
	Decomissioned            bool                                `json:"decomissioned"  validate:""`
}
