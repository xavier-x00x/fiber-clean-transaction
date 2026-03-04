package usecase

import (
	"context"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/docnumber"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"
)

type UnitUsecase struct {
	validator *validation.ValidatorHelper
	uow       transaction.UnitOfWork
	repo      repository.UnitRepository
	seqRepo   repository.NumberSequenceRepository
}

func NewUnitUsecase(r repository.UnitRepository, seqRepo repository.NumberSequenceRepository, uow transaction.UnitOfWork, v *validation.ValidatorHelper) *UnitUsecase {
	return &UnitUsecase{
		repo:      r,
		seqRepo:   seqRepo,
		uow:       uow,
		validator: v,
	}
}

func (u *UnitUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Unit, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["status"] = 1

	data, resMeta, err := u.repo.GetAllFilter(ctx, filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *UnitUsecase) FindByID(ctx context.Context, ID uint) (*entity.Unit, error) {
	data, err := u.repo.FindByID(ctx, ID)
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
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// Generate document number
		code, err := docnumber.GenerateDocNumber(ctx, u.seqRepo, "UNIT")
		if err != nil {
			return utils.Internal(err.Error(), err)
		}

		unit := &entity.Unit{
			Code:   code,
			Name:   request.Name,
			Status: request.Status,
		}

		if err := u.repo.Create(ctx, unit); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}
