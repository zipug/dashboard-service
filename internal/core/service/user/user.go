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

func (u *UserService) RegisterSupport(ctx context.Context, user models.User, created_by int64) (int64, error) {
	ok, errs := user.IsValid()
	if !ok {
		return int64(dto.BadUserId), errors.Join(errs...)
	}
	user.State = models.Registered
	dbo := dto.ToUserDbo(user)
	id, err := u.repo.RegisterSupport(ctx, dbo, created_by)
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

func (u *UserService) GetUserById(ctx context.Context, id int64) (models.User, error) {
	dbo, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	user := dbo.ToValue()
	return user, nil
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	dbos, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	var users []models.User
	for _, dbo := range dbos {
		user := dbo.ToValue()
		users = append(users, user)
	}
	return users, nil
}

func (u *UserService) GetUserEmailById(ctx context.Context, user_id int64) (models.Email, error) {
	dbo, err := u.repo.GetUserById(ctx, user_id)
	if err != nil {
		return models.Email(""), err
	}
	return models.Email(dbo.Email), nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	dbo, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	return dbo.ToValue(), nil
}

func (u *UserService) VerifyUser(ctx context.Context, user_id models.Id) error {
	err := u.repo.SetUserState(ctx, int64(user_id), models.Verified)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) SaveUser(ctx context.Context, user models.User, user_id int64) (models.User, error) {
	dbo := dto.ToUserDbo(user)
	updated_dbo, err := u.repo.SaveUser(ctx, dbo, user_id)
	if err != nil {
		return models.User{}, err
	}
	return updated_dbo.ToValue(), nil
}

func (u *UserService) ValidateUserPermissions(ctx context.Context, user_id int64, perm string) error {
	action := models.ToPermission(perm)
	return u.repo.ValidateUserPermissions(ctx, user_id, action)
}

func (u *UserService) GrantRoleToUser(ctx context.Context, user_id int64, role_id int64) error {
	return u.repo.GrantRoleToUser(ctx, user_id, role_id)
}

func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}
