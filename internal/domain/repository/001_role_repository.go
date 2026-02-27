package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type RoleRepository interface {
	GetFilter(filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error)
	FindById(id uint) (*entity.Role, error)
	FindByName(name string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, id uint, role *entity.Role) error
	AssignPermission(ctx context.Context, role *entity.Role, permissions []entity.Permission) error
	Delete(ctx context.Context, id uint) error
	ClearPermissions(ctx context.Context, role *entity.Role) error
}
