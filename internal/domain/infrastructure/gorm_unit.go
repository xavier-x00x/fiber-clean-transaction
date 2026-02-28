package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type UnitGormRepo struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) repository.UnitRepository {
	return &UnitGormRepo{db: db}
}

func (r *UnitGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Unit, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Unit{})
	return PaginateAndFilter[entity.Unit](r.db, baseQuery, filter)
}

func (r *UnitGormRepo) FindByID(ctx context.Context, ID uint) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Where("id = ?", ID).Take(&unit).Error
	return &unit, err
}

func (r *UnitGormRepo) FindByCode(ctx context.Context, code string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Where("code = ?", code).Take(&unit).Error
	return &unit, err
}

func (r *UnitGormRepo) Create(ctx context.Context, unit *entity.Unit) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(unit).Error
}

func (r *UnitGormRepo) Update(ctx context.Context, ID uint, unit *entity.Unit) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Unit{}).Where("id = ?", ID).Updates(unit).Error
}

func (r *UnitGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Unit{}, ID).Error
}
