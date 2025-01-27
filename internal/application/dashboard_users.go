package application

import (
	"context"
	"dashboard/internal/application/dto"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"fmt"
)

func (d *DashboardService) GetUserById(id int64) (models.User, error) {
	d.log.Log("info", fmt.Sprintf("getting user by id: %d", id))
	ctx := context.Background()
	usr, err := d.user.GetUserById(ctx, id)
	if err != nil {
		d.log.Log("error", fmt.Sprintf("error while getting user by id: %d", id), logger.WithErrAttr(err))
		return usr, err
	}
	d.log.Log("info", "user successfully get", logger.WithInt64Attr("user_id", id))
	return usr, nil
}

func (d *DashboardService) RegisterUser(user dto.UserDto) (int64, error) {
	d.log.Log("info", "registering user")
	ctx := context.Background()
	id, err := d.user.RegisterUser(ctx, user.ToValue())
	if err != nil {
		d.log.Log("error", "error while registering user", logger.WithErrAttr(err))
		return int64(dto.BadUserId), err
	}
	d.log.Log("info", "user successfully registered", logger.WithInt64Attr("user_id", id))
	return id, nil
}

func (d *DashboardService) LoginUser(user dto.UserDto) (models.User, error) {
	d.log.Log("info", "login user")
	ctx := context.Background()
	usr, err := d.user.LoginUser(ctx, user.ToValue())
	if err != nil {
		d.log.Log("error", "error while loginnig user", logger.WithErrAttr(err))
		return usr, err
	}
	d.log.Log("info", "user successfully logined", logger.WithInt64Attr("user_id", int64(usr.Id)))
	return usr, nil
}

func (d *DashboardService) GetAllUsers() ([]models.User, error) {
	d.log.Log("info", "get all users")
	ctx := context.Background()
	users, err := d.user.GetAllUsers(ctx)
	if err != nil {
		d.log.Log("error", "error while fetching users", logger.WithErrAttr(err))
		return users, err
	}
	d.log.Log("info", "users successfully fetched", logger.WithInt64Attr("user_count", int64(len(users))))
	return users, nil
}

func (d *DashboardService) SaveUser(user dto.UpdateUserDto) (models.User, error) {
	d.log.Log(
		"info",
		"saving user",
		logger.WithStrAttr("email", user.Email),
		logger.WithInt64Attr("id", int64(user.Id)),
	)
	ctx := context.Background()
	usr, err := d.user.SaveUser(ctx, user.ToValue())
	if err != nil {
		d.log.Log("error", "error occured while saving user", logger.WithErrAttr(err))
		return usr, err
	}
	d.log.Log("info", "user successfully saved")
	return usr, nil
}

func (d *DashboardService) DeleteUserById(user_id int64) error {
	d.log.Log("info", "deleting user", logger.WithInt64Attr("user_id", user_id))
	ctx := context.Background()
	if err := d.user.DeleteUser(ctx, user_id); err != nil {
		d.log.Log("error", "error occured while deleting user", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "user successfully deleted", logger.WithInt64Attr("user_id", user_id))
	return nil
}

func (d *DashboardService) ValidateUserPermissions(perm string, user_id int64) error {
	ctx := context.Background()
	if err := d.user.ValidateUserPermissions(ctx, user_id, perm); err != nil {
		d.log.Log("error", "error occured while validating user permissions", logger.WithErrAttr(err))
		return err
	}
	return nil
}

func (d *DashboardService) GrantRoleToUser(user_id, role_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "granting role to user", logger.WithInt64Attr("user_id", user_id), logger.WithInt64Attr("role_id", role_id))
	if err := d.user.GrantRoleToUser(ctx, user_id, role_id); err != nil {
		d.log.Log("error", "error occured while granting role to user", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "role successfully granted to user", logger.WithInt64Attr("user_id", user_id), logger.WithInt64Attr("role_id", role_id))
	return nil
}
