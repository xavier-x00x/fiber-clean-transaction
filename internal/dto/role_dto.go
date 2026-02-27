package dto

import "time"

type RoleRequest struct {
	Name string `json:"name" validate:"required,unique=roles:name"`
}

type RoleResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
