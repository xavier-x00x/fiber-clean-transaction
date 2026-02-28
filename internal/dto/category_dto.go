package dto

import "time"

type CategoryRequest struct {
	StoreCode   string `json:"store_code" validate:"required,exists=toko:code"`
	Code        string `json:"code" validate:"required,unique=categories:code"`
	Name        string `json:"name" validate:"required"`
	Status      *bool  `json:"status"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          uint       `json:"id"`
	StoreCode   string     `json:"store_code"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      *bool      `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
