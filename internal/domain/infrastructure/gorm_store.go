package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type StoreGormRepo struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) repository.StoreRepository {
	return &StoreGormRepo{db: db}
}

func (r *StoreGormRepo) GetFilter(filter entity.QueryFilter) ([]entity.Store, *entity.Meta, error) {

	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))
	searchColumn := filter.SearchColumn

	var dataList []entity.Store
	var total, totalFiltered int64

	// dataList & totalfilter
	query := r.db.Model(&entity.Store{})

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
		dataList = []entity.Store{}
	}

	// count total
	qry := r.db.Model(&entity.Store{})
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

func (r *StoreGormRepo) FindById(id uint) (*entity.Store, error) {
	var store entity.Store
	err := r.db.Where("id = ?", id).Take(&store).Error
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

func (r *StoreGormRepo) Update(ctx context.Context, id uint, store *entity.Store) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(&entity.Store{}).Where("id = ?", id).Updates(store).Error
}

func (r *StoreGormRepo) Delete(ctx context.Context, id uint) error {
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Store{}, id).Error
}
