package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type StoreGormRepo struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) repository.StoreRepository {
	return &StoreGormRepo{db: db}
}

func (r *StoreGormRepo) GetFilter(filter entity.QueryFilter) ([]entity.Store, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Store{})
	return PaginateAndFilter[entity.Store](r.db, baseQuery, filter)
}

func (r *StoreGormRepo) FindByID(ID uint) (*entity.Store, error) {
	var store entity.Store
	err := r.db.Where("id = ?", ID).Take(&store).Error
	return &store, err
}

func (r *StoreGormRepo) FindByCode(code string) (*entity.Store, error) {
	var store entity.Store
	err := r.db.Where("code = ?", code).Take(&store).Error
	return &store, err
}

func (r *StoreGormRepo) Create(ctx context.Context, store *entity.Store) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(store).Error
}

func (r *StoreGormRepo) Update(ctx context.Context, ID uint, store *entity.Store) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Store{}).Where("id = ?", ID).Updates(store).Error
}

func (r *StoreGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Store{}, ID).Error
}
