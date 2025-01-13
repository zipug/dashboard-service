package user

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"errors"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) GetUserById(ctx context.Context, id int64) (models.User, error) {
	dbo, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	user := dbo.ToValue()
	return user, nil
}

func (u *UserService) RegisterUser(ctx context.Context, user models.User) (int64, error) {
	ok, errs := user.IsValid()
	if !ok {
		return int64(dto.BadUserId), errors.Join(errs...)
	}
	user.State = models.Registered
	dbo := dto.ToUserDbo(user)
	id, err := u.repo.RegisterUser(ctx, dbo)
	if err != nil {
		return int64(dto.BadUserId), err
	}
	return id, nil
}

func (u *UserService) LoginUser(ctx context.Context, user models.User) (models.User, error) {
	ok, errs := user.IsValidForLogin()
	if !ok {
		return models.User{}, errors.Join(errs...)
	}
	user.State = models.Registered
	dbo := dto.ToUserDbo(user)
	usr, err := u.repo.GetUserByEmail(ctx, dbo.Email)
	if err != nil {
		return models.User{}, err
	}
	return usr.ToValue(), nil
}

func (u *UserService) SaveUser(ctx context.Context, user models.User) error {
	dbo := dto.ToUserDbo(user)
	return u.repo.SaveUser(ctx, dbo)
}

func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}
