package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetRoleById(role_id int64) (models.Role, error) {
	ctx := context.Background()
	d.log.Log("info", "getting role by id", logger.WithInt64Attr("role_id", role_id))
	role, err := d.role.GetRoleById(ctx, role_id)
	if err != nil {
		d.log.Log("error", "error while getting role by id", logger.WithErrAttr(err))
		return role, err
	}
	d.log.Log("info", "role successfully get", logger.WithStrAttr("role_id", role.Name))
	return role, nil
}

func (d *DashboardService) GetRoleByUserId(user_id int64) (models.Role, error) {
	ctx := context.Background()
	d.log.Log("info", "getting role by user_id", logger.WithInt64Attr("user_id", user_id))
	role, err := d.role.GetRoleByUserId(ctx, user_id)
	if err != nil {
		d.log.Log("error", "error while getting role by user_id", logger.WithErrAttr(err))
		return role, err
	}
	d.log.Log("info", "role successfully get", logger.WithStrAttr("user_id", role.Name))
	return role, nil
}

func (d *DashboardService) GetAllRoles() ([]models.Role, error) {
	ctx := context.Background()
	d.log.Log("info", "getting all roles")
	roles, err := d.role.GetAllRoles(ctx)
	if err != nil {
		d.log.Log("error", "error while fetching roles", logger.WithErrAttr(err))
		return roles, err
	}
	d.log.Log(
		"info",
		"roles successfully fetched",
		logger.WithInt64Attr("role_count", int64(len(roles))),
	)
	return roles, nil
}

func (d *DashboardService) CreateRole(role models.Role) (int64, error) {
	ctx := context.Background()
	d.log.Log("info", "creating role")
	role_id, err := d.role.CreateRole(ctx, role)
	if err != nil {
		d.log.Log("error", "error while creating role", logger.WithErrAttr(err))
		return role_id, err
	}
	d.log.Log("info", "role successfully created", logger.WithStrAttr("role_id", role.Name))
	return role_id, nil
}

func (d *DashboardService) UpdateRole(role models.Role) (models.Role, error) {
	ctx := context.Background()
	d.log.Log("info", "updating role", logger.WithInt64Attr("role_id", role.Id))
	role, err := d.role.UpdateRole(ctx, role)
	if err != nil {
		d.log.Log("error", "error while updating role", logger.WithErrAttr(err))
		return role, err
	}
	d.log.Log("info", "role successfully updated", logger.WithStrAttr("role_id", role.Name))
	return role, nil
}

func (d *DashboardService) UpdateRolePerms(user_id int64, perms []models.PermissionData) error {
	ctx := context.Background()
	d.log.Log("info", "updating role permissions", logger.WithInt64Attr("role_id", user_id))
	err := d.role.UpdateRolePermissions(ctx, user_id, perms)
	if err != nil {
		d.log.Log("error", "error while updating role permissions", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "role permissions successfully updated", logger.WithInt64Attr("role_id", user_id))
	return nil
}

func (d *DashboardService) DeleteRole(role_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "deleting role", logger.WithInt64Attr("role_id", role_id))
	err := d.role.DeleteRole(ctx, role_id)
	if err != nil {
		d.log.Log("error", "error while deleting role", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "role successfully deleted", logger.WithInt64Attr("role_id", role_id))
	return nil
}
