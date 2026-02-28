package dto

import "time"

type UnitRequest struct {
	Code   string `json:"code" validate:"required,unique=units:code"`
	Name   string `json:"name" validate:"required"`
	Status *bool  `json:"status"`
}

type UnitResponse struct {
	ID        uint       `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Status    *bool      `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
