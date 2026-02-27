package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type CategoryGormRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &CategoryGormRepo{db: db}
}

func (r *CategoryGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Category, *entity.Meta, error) {

	store_code := contextkeys.GetUserC(ctx).Store

	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))
	searchColumn := filter.SearchColumn

	var dataList []entity.Category
	var total, totalFiltered int64

	// dataList & totalfilter
	query := r.db.Model(&entity.Category{})

	// apply conditions
	query = query.Where("store_code = ?", store_code)
	if filter.Conditions != nil {
		for key, val := range filter.Conditions {
			query = query.Where(fmt.Sprintf("%s = ?", key), val)
		}
	}

	if search != "" {
		var conditions []string
		var values []interface{}

		for _, column := range searchColumn {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
			values = append(values, "%"+search+"%")
		}

		// gabungkan pakai OR
		query = query.Where("("+strings.Join(conditions, " OR ")+")", values...)
	}

	query.Count(&totalFiltered)

	query.Offset((page - 1) * limit).Limit(limit).Order(orderBy).Find(&dataList)

	if len(dataList) == 0 {
		dataList = []entity.Category{}
	}

	// count total
	qry := r.db.Model(&entity.Category{}).Where("store_code = ?", store_code)
	if filter.Conditions != nil {
		for key, val := range filter.Conditions {
			qry = qry.Where(fmt.Sprintf("%s = ?", key), val)
		}
	}
	qry.Count(&total)

	// meta
	meta := &entity.Meta{
		Page:          page,
		Limit:         limit,
		Total:         int(total),
		TotalFiltered: int(totalFiltered),
		LastPage:      int(math.Ceil(float64(totalFiltered) / float64(limit))),
		Draw:          len(dataList),
	}

	return dataList, meta, nil
}

func (r *CategoryGormRepo) FindById(ctx context.Context, id uint) (*entity.Category, error) {
	store_code := contextkeys.GetUserC(ctx).Store
	var category entity.Category
	err := r.db.Where("id = ? AND store_code = ?", id, store_code).Take(&category).Error
	return &category, err
}

func (r *CategoryGormRepo) FindByCode(ctx context.Context, code string) (*entity.Category, error) {
	store_code := contextkeys.GetUserC(ctx).Store
	var category entity.Category
	err := r.db.Where("code = ? AND store_code = ?", code, store_code).Take(&category).Error
	return &category, err
}

func (r *CategoryGormRepo) Create(ctx context.Context, category *entity.Category) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	store_code := contextkeys.GetUserC(ctx).Store
	category.StoreCode = store_code
	return gormTx.WithContext(ctx).Create(category).Error
}

func (r *CategoryGormRepo) Update(ctx context.Context, id uint, category *entity.Category) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	store_code := contextkeys.GetUserC(ctx).Store
	return gormTx.WithContext(ctx).Model(&entity.Category{}).Where("id = ? AND store_code = ?", id, store_code).Updates(category).Error
}

func (r *CategoryGormRepo) Delete(ctx context.Context, id uint) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	store_code := contextkeys.GetUserC(ctx).Store
	return gormTx.WithContext(ctx).Where("id = ? AND store_code = ?", id, store_code).Delete(&entity.Category{}).Error
}
