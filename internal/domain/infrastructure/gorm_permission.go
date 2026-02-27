package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type PermissionGormRepo struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionGormRepo {
	return &PermissionGormRepo{db: db}
}

func (r *PermissionGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Permission, *entity.Meta, error) {

	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))
	searchColumn := filter.SearchColumn

	var dataList []entity.Permission
	var total, totalFiltered int64

	// dataList & totalfilter
	query := r.db.Model(&entity.Permission{})

	// apply conditions
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
		dataList = []entity.Permission{}
	}

	// count total
	qry := r.db.Model(&entity.Permission{})
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

func (r *PermissionGormRepo) FindById(ctx context.Context, id uint) (*entity.Permission, error) {
	var data entity.Permission
	err := r.db.Where("id = ?", id).Take(&data).Error
	return &data, err
}

func (r *PermissionGormRepo) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	var data entity.Permission
	err := r.db.Where("path = ?", code).Take(&data).Error
	return &data, err
}

func (r *PermissionGormRepo) Create(ctx context.Context, permission *entity.Permission) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(permission).Error
}

func (r *PermissionGormRepo) Update(ctx context.Context, id uint, permission *entity.Permission) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Permission{}).Where("id = ?", id).Updates(permission).Error
}

func (r *PermissionGormRepo) Delete(ctx context.Context, id uint) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Permission{}, id).Error
}
