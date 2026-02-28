package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"

	"gorm.io/gorm"
)

type CategoryGormRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &CategoryGormRepo{db: db}
}

func (r *CategoryGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Category, *entity.Meta, error) {
	storeCode := contextkeys.GetUserC(ctx).Store
	baseQuery := r.db.Model(&entity.Category{}).Where("store_code = ?", storeCode)
	return PaginateAndFilter[entity.Category](r.db, baseQuery, filter)
}

func (r *CategoryGormRepo) FindByID(ctx context.Context, ID uint) (*entity.Category, error) {
	storeCode := contextkeys.GetUserC(ctx).Store
	var category entity.Category
	err := r.db.Where("id = ? AND store_code = ?", ID, storeCode).Take(&category).Error
	return &category, err
}

func (r *CategoryGormRepo) FindByCode(ctx context.Context, code string) (*entity.Category, error) {
	storeCode := contextkeys.GetUserC(ctx).Store
	var category entity.Category
	err := r.db.Where("code = ? AND store_code = ?", code, storeCode).Take(&category).Error
	return &category, err
}

func (r *CategoryGormRepo) Create(ctx context.Context, category *entity.Category) error {
	gormTx := GetDBWithTx(ctx, r.db)
	storeCode := contextkeys.GetUserC(ctx).Store
	category.StoreCode = storeCode
	return gormTx.WithContext(ctx).Create(category).Error
}

func (r *CategoryGormRepo) Update(ctx context.Context, ID uint, category *entity.Category) error {
	gormTx := GetDBWithTx(ctx, r.db)
	storeCode := contextkeys.GetUserC(ctx).Store
	return gormTx.WithContext(ctx).Model(&entity.Category{}).Where("id = ? AND store_code = ?", ID, storeCode).Updates(category).Error
}

func (r *CategoryGormRepo) Delete(ctx context.Context, ID uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	storeCode := contextkeys.GetUserC(ctx).Store
	return gormTx.WithContext(ctx).Where("id = ? AND store_code = ?", ID, storeCode).Delete(&entity.Category{}).Error
}
