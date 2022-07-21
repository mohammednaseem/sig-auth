package model

type JWT struct {
	Key       string `json:"key" validate:"required"`
	Payload   string `json:"payload" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}
