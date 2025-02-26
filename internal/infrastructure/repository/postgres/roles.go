package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrGetRoleById     = errors.New("could not get role by id")
	ErrGetRoleByUserId = errors.New("could not get role by user id")
	ErrGetRoles        = errors.New("could not get roles")
	ErrCreateRole      = errors.New("could not create role")
	ErrUpdateRole      = errors.New("could not update role")
	ErrUpdatePerms     = errors.New("could not update role permissions")
	ErrDeleteRole      = errors.New("could not delete role")
)

func (repo *PostgresRepository) GetRoleById(ctx context.Context, role_id int64) (*dto.RolesDbo, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		SELECT r.id, r.name, r.description, r.is_custom
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
		SELECT p.name, p.description, rp.do_create, rp.do_read, rp.do_update, rp.do_delete
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

func (repo *PostgresRepository) GetRoleByUserId(ctx context.Context, user_id int64) (*dto.RolesDbo, error) {
	role_user_rows, err := pu.Dispatch[dto.UserRolesDbo](
		ctx,
		repo.db,
		`
		SELECT user_id, role_id
		FROM user_roles
		WHERE user_id = $1::bigint;
		`,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(role_user_rows) == 0 {
		return nil, ErrGetRoleByUserId
	}
	return repo.GetRoleById(ctx, role_user_rows[0].RoleId)
}

func (repo *PostgresRepository) GetAllRoles(ctx context.Context) ([]dto.RolesDbo, error) {
	role_rows, err := pu.Dispatch[dto.RolesDbo](
		ctx,
		repo.db,
		`
		SELECT r.id, r.name, r.description, r.is_custom
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
		SELECT rp.role_id, p.name, p.description, rp.do_create, rp.do_read, rp.do_update, rp.do_delete
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
	for i := range role_rows {
		for j := range perms_rows {
			if role_rows[i].Id == perms_rows[j].RoleId {
				role_rows[i].Permissions = append(role_rows[i].Permissions, perms_rows[j])
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
		INSERT INTO roles (name, description, is_custom)
		VALUES ($1::text, $2::text, true)
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
		  AND is_custom is FALSE
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
	var errs []error
	tx := repo.db.MustBegin()
UpdateLoop:
	for _, perm := range perms {
		_, err := pu.DispatchTx[dto.RolePermissionDbo](
			ctx,
			tx,
			`
			UPDATE role_permissions rp
			SET do_create = $1::bool,
					do_read = $2::bool,
					do_update = $3::bool,
					do_delete = $4::bool
			FROM permissions p
				WHERE p.name = $5::text
					AND rp.role_id = $6::bigint
					AND rp.permission_id = p.id;
			`,
			perm.DoCreate,
			perm.DoRead,
			perm.DoUpdate,
			perm.DoDelete,
			perm.Name,
			role_id,
		)
		if err != nil {
			errs = append(errs, err)
			tx.Rollback()
			break UpdateLoop
		}
	}

	if len(errs) > 0 {
		return ErrUpdatePerms
	}
	tx.Commit()

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
			)
			AND is_custom is FALSE;
		`,
		role_id,
	)
	if err != nil {
		return err
	}
	return nil
}
