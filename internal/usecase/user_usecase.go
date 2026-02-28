package usecase

import (
	"errors"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/domain/repository"
	"fiber-clean-transaction/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: r,
	}
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

func (u *UserUsecase) UpdateAvatar(ID uint, avatar string) error {
	return u.repo.UpdateAvatar(ID, avatar)
}
