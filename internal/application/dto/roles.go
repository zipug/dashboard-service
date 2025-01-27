package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type PermissionDto struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	DoCreate bool   `json:"do_create,omitempty"`
	DoUpdate bool   `json:"do_update,omitempty"`
	DoRead   bool   `json:"do_read,omitempty"`
	DoDelete bool   `json:"dodelete,omitempty"`
}

type RoleDto struct {
	Id          int64           `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Permissions []PermissionDto `json:"permissions,omitempty"`
}

type RolesDbo struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Permissions []RolePermissionDbo
	CreatedAt   string         `db:"created_at,omitempty"`
	UpdateAt    string         `db:"updated_at,omitempty"`
	DeleteAt    sql.NullString `db:"deleted_at,omitempty"`
}

type RolePermissionDbo struct {
	RoleId       int64          `db:"role_id,omitempty" json:"role_id,omitempty"`
	PermissionId int64          `db:"permission_id,omitempty" json:"permission_id,omitempty"`
	Name         sql.NullString `db:"name,omitempty" json:"name,omitempty"`
	DoCreate     bool           `db:"do_create" json:"do_create"`
	DoRead       bool           `db:"do_read" json:"do_read"`
	DoUpdate     bool           `db:"do_update" json:"do_update"`
	DoDelete     bool           `db:"do_delete" json:"do_delete"`
}

func ToRoleDto(r models.Role) RoleDto {
	return RoleDto{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Permissions: ToPermissionDto(r.Permissions),
	}
}

func ToPermissionDto(p []models.PermissionData) []PermissionDto {
	result := make([]PermissionDto, 0)
	for _, perm := range p {
		result = append(result, PermissionDto{
			Name:     perm.Name,
			DoCreate: perm.Create,
			DoRead:   perm.Read,
			DoUpdate: perm.Update,
			DoDelete: perm.Delete,
		})
	}
	return result
}

func ToRoleDbo(r models.Role) RolesDbo {
	return RolesDbo{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
		Permissions: ToRolePermissionDbo(r.Permissions),
	}
}

func ToRolePermissionDbo(p []models.PermissionData) []RolePermissionDbo {
	result := make([]RolePermissionDbo, 0)
	for _, perm := range p {
		result = append(result, RolePermissionDbo{
			Name:     sql.NullString{String: perm.Name, Valid: true},
			DoCreate: perm.Create,
			DoRead:   perm.Read,
			DoUpdate: perm.Update,
			DoDelete: perm.Delete,
		})
	}
	return result
}

func (pd *PermissionDto) ToValue() models.PermissionData {
	return models.PermissionData{
		Id:     pd.Id,
		Name:   pd.Name,
		Create: pd.DoCreate,
		Read:   pd.DoRead,
		Update: pd.DoUpdate,
		Delete: pd.DoDelete,
	}
}

func (pd *RolePermissionDbo) ToValue() models.PermissionData {
	return models.PermissionData{
		Id:     pd.PermissionId,
		Name:   pd.Name.String,
		Create: pd.DoCreate,
		Read:   pd.DoRead,
		Update: pd.DoUpdate,
		Delete: pd.DoDelete,
	}
}

func (rd *RoleDto) ToValue() models.Role {
	var perms []models.PermissionData
	for _, perm := range rd.Permissions {
		perms = append(perms, models.PermissionData{
			Name:   perm.Name,
			Create: perm.DoCreate,
			Read:   perm.DoRead,
			Update: perm.DoUpdate,
			Delete: perm.DoDelete,
		})
	}
	return models.Role{
		Id:          rd.Id,
		Name:        rd.Name,
		Description: rd.Description,
		Permissions: perms,
	}
}

func (rd *RolesDbo) ToValue() models.Role {
	var perms []models.PermissionData
	for _, perm := range rd.Permissions {
		perms = append(perms, models.PermissionData{
			Name:   perm.Name.String,
			Create: perm.DoCreate,
			Read:   perm.DoRead,
			Update: perm.DoUpdate,
			Delete: perm.DoDelete,
		})
	}
	return models.Role{
		Id:          rd.Id,
		Name:        rd.Name,
		Description: rd.Description,
		Permissions: perms,
	}
}

func (rp *RolePermissionDbo) Can(action models.PermissionAction) bool {
	switch action {
	case models.Create:
		if rp.DoCreate {
			return true
		}
		return false
	case models.Update:
		if rp.DoUpdate {
			return true
		}
		return false
	case models.Read:
		if rp.DoRead {
			return true
		}
		return false
	case models.Delete:
		if rp.DoDelete {
			return true
		}
		return false
	default:
		return false
	}
}
