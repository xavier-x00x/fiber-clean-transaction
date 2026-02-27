package repository

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
)

type UnitRepository interface {
	GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Unit, *entity.Meta, error)
	FindById(ctx context.Context, id uint) (*entity.Unit, error)
	FindByCode(ctx context.Context, code string) (*entity.Unit, error)
	Create(ctx context.Context, unit *entity.Unit) error
	Update(ctx context.Context, id uint, unit *entity.Unit) error
	Delete(ctx context.Context, id uint) error
}
