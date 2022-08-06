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
