package usecase

import (
	"context"
	"errors"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"
	trans "fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	validator *validation.ValidatorHelper
	uow       trans.UnitOfWork
	repo      repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository, v *validation.ValidatorHelper, u trans.UnitOfWork) *UserUsecase {
	return &UserUsecase{
		repo:      r,
		validator: v,
		uow:       u,
	}
}

func (u *UserUsecase) GetAllFilter(meta *dto.MetaRequest) ([]entity.User, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "username", "email", "updated_at"}
	searchColumns := []string{"id", "name", "username", "email"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.repo.GetFilter(filter)
	if err != nil {
		return nil, nil, utils.Internal(err.Error(), err)
	}
	return data, resMeta, nil
}

func (u *UserUsecase) Register(request *dto.UserRequest) error {

	data := &entity.User{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	data.Password = string(hashed)
	return u.repo.Create(data)
}

func (u *UserUsecase) Login(email, password string) (*entity.User, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("password salah")
	}

	return user, nil
}

func (u *UserUsecase) Profile(ID uint) (*entity.User, error) {
	return u.repo.FindByID(ID)
}

func (u *UserUsecase) GoogleProfile(email string) (*entity.User, error) {
	return u.repo.FindByEmail(email)
}

func (u *UserUsecase) Update(ctx context.Context, ID uint, request *dto.UserUpdateRequest) error {
	// Validasi SEBELUM transaksi
	if err := u.validator.ValidateUpdate(request, ID); err != nil {
		return err
	}

	return u.uow.Do(ctx, func(ctx context.Context) error {
		if _, err := u.repo.FindByID(ID); err != nil {
			return utils.NotFound(err.Error())
		}

		data := &entity.User{
			Name:     request.Name,
			Username: request.Username,
			Email:    request.Email,
			Role:     request.Role,
		}

		if request.Password != "" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			data.Password = string(hashed)
		}

		if err := u.repo.Update(ctx, ID, data); err != nil {
			return utils.Internal(err.Error(), err)
		}

		return nil
	})
}

func (u *UserUsecase) UpdateAvatar(ID uint, avatar string) error {
	return u.repo.UpdateAvatar(ID, avatar)
}

func (u *UserUsecase) Delete(ctx context.Context, ID uint) error {
	return u.uow.Do(ctx, func(ctx context.Context) error {
		// cek apakah data ada
		if _, err := u.repo.FindByID(ID); err != nil {
			return utils.NotFound(err.Error())
		}
		// hapus data
		if err := u.repo.Delete(ctx, ID); err != nil {
			return utils.Internal(err.Error(), err)
		}
		return nil
	})
}
