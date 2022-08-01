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
