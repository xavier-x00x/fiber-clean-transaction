package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type PermissionGormRepo struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionGormRepo{db: db}
}

func (r *PermissionGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Permission, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Permission{})
	return PaginateAndFilter[entity.Permission](r.db, baseQuery, filter)
}

func (r *PermissionGormRepo) FindByID(ctx context.Context, ID uint) (*entity.Permission, error) {
	var data entity.Permission
	err := r.db.Where("id = ?", ID).Take(&data).Error
	return &data, err
}

func (r *PermissionGormRepo) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	var data entity.Permission
	err := r.db.Where("path = ?", code).Take(&data).Error
	return &data, err
}

func (r *PermissionGormRepo) Create(ctx context.Context, permission *entity.Permission) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(permission).Error
}

func (r *PermissionGormRepo) Update(ctx context.Context, ID uint, permission *entity.Permission) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Permission{}).Where("id = ?", ID).Updates(permission).Error
}

func (r *PermissionGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Permission{}, ID).Error
}
