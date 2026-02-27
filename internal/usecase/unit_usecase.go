package usecase

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"
)

type UnitUsecase struct {
	repo      repository.UnitRepository
	txManager transaction.TransactionManager
	validator *validation.ValidatorHelper
}

func NewUnitUsecase(r repository.UnitRepository, t transaction.TransactionManager, v *validation.ValidatorHelper) *UnitUsecase {
	return &UnitUsecase{
		repo:      r,
		txManager: t,
		validator: v,
	}
}

func (u *UnitUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Unit, *entity.Meta, error) {
	// validasi order by untuk hindari SQL injection
	direction := []string{"asc", "desc"}
	order := []string{"id", "code", "name", "updated_at"}
	searchColumn := []string{"id", "code", "name"}

	if !utils.Contains(order, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = "id"
		meta.OrderDir = "asc"
	}

	filter := entity.QueryFilter{
		Page:         meta.Page,
		Limit:        meta.Limit,
		Search:       meta.Search,
		OrderColumn:  meta.OrderColumn,
		OrderDir:     meta.OrderDir,
		SearchColumn: searchColumn,
		Conditions:   map[string]interface{}{},
	}

	filter.Conditions["status"] = 1

	data, resMeta, err := u.repo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *UnitUsecase) FindById(ctx context.Context, id uint) (*entity.Unit, error) {
	data, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *UnitUsecase) FindByCode(ctx context.Context, code string) (*entity.Unit, error) {
	data, err := u.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *UnitUsecase) Create(ctx context.Context, request *dto.UnitRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	txErr := u.txManager.WithTransaction(ctx, func(ctx context.Context, tx transaction.Transaction) error {

		unit := &entity.Unit{
			Code:   request.Code,
			Name:   request.Name,
			Status: request.Status,
		}

		if err := u.repo.Create(ctx, tx, unit); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})

	return txErr
}
