package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type RolesService interface {
	GetRoleById(context.Context, int64) (models.Role, error)
	GetRoleByUserId(context.Context, int64) (models.Role, error)
	GetAllRoles(context.Context) ([]models.Role, error)
	CreateRole(context.Context, models.Role) (int64, error)
	UpdateRole(context.Context, models.Role) (models.Role, error)
	UpdateRolePermissions(context.Context, int64, []models.PermissionData) error
	DeleteRole(context.Context, int64) error
}

type RoleRepository interface {
	GetRoleById(context.Context, int64) (*dto.RolesDbo, error)
	GetRoleByUserId(context.Context, int64) (*dto.RolesDbo, error)
	GetAllRoles(context.Context) ([]dto.RolesDbo, error)
	CreateRole(context.Context, dto.RolesDbo) (int64, error)
	UpdateRole(context.Context, dto.RolesDbo) (*dto.RolesDbo, error)
	UpdateRolePermissions(context.Context, int64, []dto.RolePermissionDbo) error
	DeleteRole(context.Context, int64) error
}
