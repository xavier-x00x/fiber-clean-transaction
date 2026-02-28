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

type StoreUsecase struct {
	validator *validation.ValidatorHelper
	uow       transaction.UnitOfWork
	repo      repository.StoreRepository
}

func NewStoreUsecase(r repository.StoreRepository, uow transaction.UnitOfWork, v *validation.ValidatorHelper) *StoreUsecase {
	return &StoreUsecase{
		repo:      r,
		uow:       uow,
		validator: v,
	}
}

func (u *StoreUsecase) GetAllFilter(meta *dto.MetaRequest) ([]entity.Store, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["status"] = 1

	data, resMeta, err := u.repo.GetFilter(filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *StoreUsecase) FindByID(ID uint) (*entity.Store, error) {
	data, err := u.repo.FindByID(ID)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *StoreUsecase) FindByCode(code string) (*entity.Store, error) {
	data, err := u.repo.FindByCode(code)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *StoreUsecase) Create(ctx context.Context, request *dto.StoreRequest) error {

	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}

	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		store := &entity.Store{
			Code:    request.Code,
			Name:    request.Name,
			Npwp:    request.Npwp,
			Address: request.Address,
			Phone:   request.Phone,
			Email:   request.Email,
			Phone2:  request.Phone2,
			Email2:  request.Email2,
			Status:  request.Status,
		}

		if err := u.repo.Create(ctx, store); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *StoreUsecase) Update(ctx context.Context, ID uint, request *dto.StoreRequest) error {

	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, ID); err != nil {
		return err
	}

	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindByID(ID); err != nil {
			return utils.NotFound(err.Error())
		}

		store := &entity.Store{
			Code:    request.Code,
			Name:    request.Name,
			Npwp:    request.Npwp,
			Address: request.Address,
			Phone:   request.Phone,
			Email:   request.Email,
			Phone2:  request.Phone2,
			Email2:  request.Email2,
			Status:  request.Status,
		}

		if err := u.repo.Update(ctx, ID, store); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *StoreUsecase) Delete(ctx context.Context, ID uint) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindByID(ID); err != nil {
			return utils.NotFound(err.Error())
		}

		// hapus data
		if err := u.repo.Delete(ctx, ID); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil // berhasil di commit
	})
}
