package dto

import "time"

type PermissionRequest struct {
	Path string `json:"path" validate:"required,unique=permissions:path"`
	Name string `json:"name" validate:"required,unique=permissions:name"`
}

type PermissionResponse struct {
	ID        uint       `json:"id"`
	Path      string     `json:"path"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
