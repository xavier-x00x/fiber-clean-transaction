package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type RoleRepository interface {
	GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error)
	FindByID(ctx context.Context, ID uint) (*entity.Role, error)
	FindByName(name string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, ID uint, role *entity.Role) error
	AssignPermission(ctx context.Context, role *entity.Role, permissions []entity.Permission) error
	Delete(ctx context.Context, ID uint) error
	ClearPermissions(ctx context.Context, role *entity.Role) error
}
