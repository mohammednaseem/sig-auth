package model

type Device struct {
	DeviceId    string   `json:"deviceid" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Password    string   `json:"password"`
	Cerificate1 string   `json:"cerificate1"`
	Cerificate2 string   `json:"cerificate2"`
	Cerificate3 string   `json:"cerificate3"`
	Project     string   `json:"Project"`
	Region      string   `json:"Region"`
	Created_On  string   `json:"created_on"`
	PublicKey   []string `json:"key" validate:"required"`
}
