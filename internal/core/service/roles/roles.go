package roles

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"database/sql"
)

type RolesService struct {
	repo ports.RoleRepository
}

func NewRolesService(repo ports.RoleRepository) *RolesService {
	return &RolesService{repo: repo}
}

func (r *RolesService) GetRoleById(ctx context.Context, user_id int64) (models.Role, error) {
	roleDbo, err := r.repo.GetRoleById(ctx, user_id)
	if err != nil {
		return models.Role{}, err
	}
	role := roleDbo.ToValue()
	return role, nil
}

func (r *RolesService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	rolesDbo, err := r.repo.GetAllRoles(ctx)
	if err != nil {
		return []models.Role{}, err
	}
	var roles []models.Role
	for _, role := range rolesDbo {
		roles = append(roles, role.ToValue())
	}
	return roles, nil
}

func (r *RolesService) CreateRole(ctx context.Context, role models.Role) (int64, error) {
	roleDbo := dto.ToRoleDbo(role)
	role_id, err := r.repo.CreateRole(ctx, roleDbo)
	if err != nil {
		return role_id, err
	}
	return role_id, nil
}

func (r *RolesService) UpdateRole(ctx context.Context, role models.Role) (models.Role, error) {
	roleDbo := dto.ToRoleDbo(role)
	newRole, err := r.repo.UpdateRole(ctx, roleDbo)
	if err != nil {
		return models.Role{}, err
	}
	return newRole.ToValue(), nil
}

func (r *RolesService) UpdateRolePermissions(ctx context.Context, role_id int64, perms []models.PermissionData) error {
	var permsDbo []dto.RolePermissionDbo
	for _, perm := range perms {
		permsDbo = append(permsDbo, dto.RolePermissionDbo{
			PermissionId: perm.Id,
			Name:         sql.NullString{String: perm.Name, Valid: true},
			DoCreate:     perm.Create,
			DoRead:       perm.Read,
			DoUpdate:     perm.Update,
			DoDelete:     perm.Delete,
		})
	}
	return r.repo.UpdateRolePermissions(ctx, role_id, permsDbo)
}

func (r *RolesService) DeleteRole(ctx context.Context, user_id int64) error {
	return r.repo.DeleteRole(ctx, user_id)
}
