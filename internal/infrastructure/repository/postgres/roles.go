package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrGetRoleById = errors.New("could not get role by id")
	ErrGetRoles    = errors.New("could not get roles")
	ErrCreateRole  = errors.New("could not create role")
	ErrUpdateRole  = errors.New("could not update role")
	ErrUpdatePerms = errors.New("could not update role permissions")
	ErrDeleteRole  = errors.New("could not delete role")
)

func (repo *PostgresRepository) GetRoleById(ctx context.Context, role_id int64) (*dto.RolesDbo, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		SELECT r.id, r.name, r.description
		FROM roles r
		WHERE r.id = $1::bigint;
		`,
		role_id,
	)
	if err != nil {
		return nil, err
	}
	if len(role_rows) == 0 {
		return nil, ErrGetRoleById
	}
	perms_rows, err := pu.Dispatch[dto.RolePermissionDbo](
		ctx,
		repo.db,
		`
		SELECT p.name, rp.do_create, rp.do_read, rp.do_update, rp.do_delete
		FROM role_permissions rp
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = $1::bigint;
		`,
		role_id,
	)
	if err != nil {
		return nil, err
	}
	if len(perms_rows) == 0 {
		return nil, ErrGetRoleById
	}
	role := role_rows[0]
	role.Permissions = perms_rows
	return &role, nil
}

func (repo *PostgresRepository) GetAllRoles(ctx context.Context) ([]dto.RolesDbo, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		SELECT r.id, r.name, r.description
		FROM roles r
		WHERE r.deleted_at IS NULL;
		`,
	)
	if err != nil {
		return nil, err
	}
	if len(role_rows) == 0 {
		return nil, ErrGetRoles
	}
	perms_rows, err := pu.Dispatch[dto.RolePermissionDbo](
		ctx,
		repo.db,
		`
		SELECT rp.role_id, p.name, rp.do_create, rp.do_read, rp.do_update, rp.do_delete
		FROM role_permissions rp
		LEFT JOIN permissions p ON rp.permission_id = p.id
		`,
	)
	if err != nil {
		return nil, err
	}
	if len(perms_rows) == 0 {
		return nil, ErrGetRoles
	}
	for _, role := range role_rows {
		for _, perm := range perms_rows {
			if role.Id == perm.RoleId {
				role.Permissions = append(role.Permissions, perm)
			}
		}
	}
	return role_rows, nil
}

func (repo *PostgresRepository) CreateRole(ctx context.Context, role dto.RolesDbo) (int64, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		INSERT INTO roles (name, description)
		VALUES ($1::text, $2::text)
		RETURNING *;
		`,
		role.Name,
		role.Description,
	)
	if err != nil {
		return -1, err
	}
	if len(role_rows) == 0 {
		return -1, ErrCreateRole
	}
	new_role := role_rows[0]
	perms_rows, err := pu.Dispatch[dto.RolePermissionDbo](
		ctx,
		repo.db,
		`
		INSERT INTO role_permissions (role_id, permission_id, do_create, do_read, do_update, do_delete)
		SELECT $1::bigint, p.id, false, false, false, false
		FROM permissions p
		RETURNING *;
		`,
		new_role.Id,
	)
	if err != nil {
		return -1, err
	}
	if len(perms_rows) == 0 {
		return -1, ErrGetRoleById
	}
	return new_role.Id, nil
}

func (repo *PostgresRepository) UpdateRole(ctx context.Context, role dto.RolesDbo) (*dto.RolesDbo, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		UPDATE roles
		SET name = $1::text, description = $2::text
		WHERE id = $3::bigint
		RETURNING *;
		`,
		role.Name,
		role.Description,
		role.Id,
	)
	if err != nil {
		return nil, err
	}
	if len(role_rows) == 0 {
		return nil, ErrUpdateRole
	}
	new_role := role_rows[0]
	return &new_role, nil
}

func (repo *PostgresRepository) UpdateRolePermissions(ctx context.Context, role_id int64, perms []dto.RolePermissionDbo) error {
	_, err := pu.Dispatch[dto.RolePermissionDbo](
		ctx,
		repo.db,
		`
		UPDATE role_permissions rp
		SET permission_id = t.permission_id,
		do_create = t.do_create,
		do_read = t.do_read,
		do_update = t.do_update,
		do_delete = t.do_delete
		FROM (
			--)
			VALUES(:permission_id, :do_create, :do_read, :do_update, :do_delete)
		) AS t (permission_id, do_create, do_read, do_update, do_delete)
		WHERE rp.role_id = $1::bigint;
		`,
		role_id,
		perms,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) DeleteRole(ctx context.Context, role_id int64) error {
	_, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		DELETE
		FROM roles
		WHERE id = $1::bigint
		  AND id IN (
      	SELECT r.id
				FROM roles r
				WHERE r.id = $1::bigint
				  AND r.name <> 'admin'
			);
		`,
		role_id,
	)
	if err != nil {
		return err
	}
	return nil
}
