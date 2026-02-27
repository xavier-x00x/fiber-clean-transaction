package infrastructure

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type RoleGormRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleGormRepo {
	return &RoleGormRepo{db: db}
}

func (r *RoleGormRepo) GetAllFilter(ctx context.Context, filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error) {
	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))
	searchColumn := filter.SearchColumn

	var dataList []entity.Role
	var total, totalFiltered int64

	// dataList & totalfilter
	query := r.db.Model(&entity.Role{})

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
		dataList = []entity.Role{}
	}

	// count total
	qry := r.db.Model(&entity.Role{})
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

func (r *RoleGormRepo) FindById(ctx context.Context, id uint) (*entity.Role, error) {
	var data entity.Role
	err := r.db.Where("id = ?", id).Take(&data).Error
	return &data, err
}

func (r *RoleGormRepo) FindByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", name).Take(&role).Error
	return &role, err
}

func (r *RoleGormRepo) Create(ctx context.Context, data *entity.Role) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Create(data).Error
}

func (r *RoleGormRepo) Update(ctx context.Context, data *entity.Role) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Save(data).Error
}

func (r *RoleGormRepo) AssignPermission(ctx context.Context, role *entity.Role, permissions []entity.Permission) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

func (r *RoleGormRepo) Delete(ctx context.Context, id uint) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Delete(&entity.Role{}, id).Error
}

func (r *RoleGormRepo) ClearPermissions(ctx context.Context, role *entity.Role) error {
	// Type assertion untuk mengkonversi abstraksi ke implementasi konkret
	gormTx := GetDBWithTx(ctx, r.db)
	return gormTx.WithContext(ctx).Model(role).Association("Permissions").Clear()
}
