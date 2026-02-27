package dto

import "time"

type StoreRequest struct {
	Code    string `json:"code" validate:"required,unique=stores:code"`
	Name    string `json:"name" validate:"required"`
	Npwp    string `json:"npwp"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email" validate:"required,email"`
	Phone2  string `json:"phone2"`
	Email2  string `json:"email2"`
	Status  bool   `json:"status"`
}

type StoreResponse struct {
	ID        uint       `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Npwp      string     `json:"npwp"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Phone2    string     `json:"phone2"`
	Email2    string     `json:"email2"`
	Status    bool       `json:"status"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (t StoreRequest) TableName() string {
	return "stores"
}

func (t StoreResponse) TableName() string {
	return "stores"
}
