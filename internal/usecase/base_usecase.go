package usecase

import (
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/pkg/utils"
)

// BuildQueryFilter creates a standardized QueryFilter from a MetaRequest,
// validating the order column/direction against allowed values to prevent SQL injection.
func BuildQueryFilter(meta *dto.MetaRequest, allowedOrder []string, searchColumns []string) entity.QueryFilter {
	direction := []string{"asc", "desc"}

	if !utils.Contains(allowedOrder, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = "id"
		meta.OrderDir = "asc"
	}

	return entity.QueryFilter{
		Page:         meta.Page,
		Limit:        meta.Limit,
		Search:       meta.Search,
		OrderColumn:  meta.OrderColumn,
		OrderDir:     meta.OrderDir,
		SearchColumn: searchColumns,
		Conditions:   map[string]interface{}{},
	}
}
