package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/transaction"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type UnitGormRepo struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) repository.UnitRepository {
	return &UnitGormRepo{db: db}
}

func (r *UnitGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Unit, *entity.Meta, error) {

	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))
	searchColumn := filter.SearchColumn

	var dataList []entity.Unit
	var total, totalFiltered int64

	// dataList & totalfilter
	query := r.db.Model(&entity.Unit{})

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
		dataList = []entity.Unit{}
	}

	// count total
	qry := r.db.Model(&entity.Unit{})
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

func (r *UnitGormRepo) FindById(ctx context.Context, id uint) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Where("id = ?", id).Take(&unit).Error
	return &unit, err
}

func (r *UnitGormRepo) FindByCode(ctx context.Context, code string) (*entity.Unit, error) {
	var unit entity.Unit
	err := r.db.Where("code = ?", code).Take(&unit).Error
	return &unit, err
}

func (r *UnitGormRepo) Create(ctx context.Context, tx transaction.Transaction, unit *entity.Unit) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(unit).Error
}

func (r *UnitGormRepo) Update(ctx context.Context, tx transaction.Transaction, id uint, unit *entity.Unit) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Unit{}).Where("id = ?", id).Updates(unit).Error
}

func (r *UnitGormRepo) Delete(ctx context.Context, tx transaction.Transaction, id uint) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Unit{}, id).Error
}
