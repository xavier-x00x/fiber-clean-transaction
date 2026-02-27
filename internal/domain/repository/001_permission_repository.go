package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type PermissionRepository interface {
	GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Permission, *entity.Meta, error)
	FindById(ctx context.Context, id uint) (*entity.Permission, error)
	FindByCode(ctx context.Context, code string) (*entity.Permission, error)
	Create(ctx context.Context, permission *entity.Permission) error
	Update(ctx context.Context, id uint, permission *entity.Permission) error
	Delete(ctx context.Context, id uint) error
}
