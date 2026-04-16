package dto

import "time"

type RoleRequest struct {
	Name        string   `json:"name" validate:"required,unique=roles:name"`
	Permissions []string `json:"permissions" validate:"required,dive"`
}

type RoleResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Permissions []string   `json:"permissions"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
