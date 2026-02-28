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

type CategoryUsecase struct {
	validator *validation.ValidatorHelper
	uow       transaction.UnitOfWork
	repo      repository.CategoryRepository
}

func NewCategoryUsecase(r repository.CategoryRepository, uow transaction.UnitOfWork, v *validation.ValidatorHelper) *CategoryUsecase {
	return &CategoryUsecase{
		validator: v,
		uow:       uow,
		repo:      r,
	}
}

func (u *CategoryUsecase) GetAllFilter(ctx context.Context, meta *dto.MetaRequest) ([]entity.Category, *entity.Meta, error) {
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

func (u *CategoryUsecase) FindByID(ctx context.Context, ID uint) (*entity.Category, error) {
	data, err := u.repo.FindByID(ctx, ID)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *CategoryUsecase) FindByCode(ctx context.Context, code string) (*entity.Category, error) {
	data, err := u.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, utils.NotFound(err.Error())
	}
	return data, nil
}

func (u *CategoryUsecase) Create(ctx context.Context, request *dto.CategoryRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.Validate(request); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		category := &entity.Category{
			StoreCode:   request.StoreCode,
			Code:        request.Code,
			Name:        request.Name,
			Description: request.Description,
			Status:      request.Status,
		}

		if err := u.repo.Create(ctx, category); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *CategoryUsecase) Update(ctx context.Context, ID uint, request *dto.CategoryRequest) error {
	// ✅ Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, ID); err != nil {
		return err
	}
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		category := &entity.Category{
			StoreCode:   request.StoreCode,
			Code:        request.Code,
			Name:        request.Name,
			Description: request.Description,
			Status:      request.Status,
		}

		if err := u.repo.Update(ctx, ID, category); err != nil {
			return utils.Internal(err.Error(), err)
		}
		// lakukan logika lain ...

		return nil // berhasil di commit
	})
}

func (u *CategoryUsecase) Delete(ctx context.Context, ID uint) error {
	// ✅ UseCase tidak tahu tentang GORM, hanya menggunakan abstraksi
	return u.uow.Do(ctx, func(ctx context.Context) error {

		// cek apakah data ada
		if _, err := u.repo.FindByID(ctx, ID); err != nil {
			return utils.NotFound(err.Error())
		}

		// hapus data
		if err := u.repo.Delete(ctx, ID); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil // berhasil di commit
	})
}
