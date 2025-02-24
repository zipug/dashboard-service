package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	pu "dashboard/pkg/postgres_utils"
	"errors"
)

var (
	ErrRegister             = errors.New("could not register user")
	ErrUpdate               = errors.New("could not update user")
	ErrPqEmailAlreadyExists = errors.New("pq: duplicate key value violates unique constraint \"users_email_unique\"")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrGetUserById          = errors.New("could not get user by id")
	ErrGetUsers             = errors.New("could not get users")
	ErrValidatePermissions  = errors.New("could not validate permissions")
)

func (repo *PostgresRepository) RegisterUser(ctx context.Context, user dto.UserDbo) (int64, error) {
	tx := repo.db.MustBegin()
	usr_rows, err := pu.DispatchTx[dto.UserDbo](
		ctx,
		tx,
		`
		INSERT INTO users (state, email, password, name, lastname, avatar_url)
		VALUES($1::text, $2::text, $3::text, $4::text, $5::text, $6::text)
		RETURNING *;
		`,
		user.State,
		user.Email,
		user.Password,
		user.Name,
		user.LastName,
		user.AvatarUrl,
	)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, ErrPqEmailAlreadyExists) {
			return int64(dto.BadUserId), ErrEmailAlreadyExists
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_unique\"" {
			return int64(dto.BadUserId), ErrEmailAlreadyExists
		}
		return int64(dto.BadUserId), err
	}
	if len(usr_rows) == 0 {
		tx.Rollback()
		return int64(dto.BadUserId), ErrRegister
	}
	usr := usr_rows[0]
	_, err = tx.Queryx(
		`
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1::bigint AS user_id, r.id AS role_id
		FROM roles r
		WHERE r.name = $2::text
		`,
		usr.Id,
		"user",
	)
	if err != nil {
		tx.Rollback()
		return int64(dto.BadUserId), err
	}
	tx.Commit()
	return usr.Id, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*dto.UserDbo, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		SELECT id, state, email, password, name, lastname, avatar_url
		FROM users
		WHERE email = $1::text;
		`,
		email,
	)
	if err != nil {
		return nil, err
	}
	if len(usr_rows) == 0 {
		return nil, ErrRegister
	}
	usr := usr_rows[0]
	return &usr, nil
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id int64) (*dto.UserDbo, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		SELECT id, state, email, name, lastname, avatar_url
		FROM users
		WHERE id = $1::bigint;
		`,
		id,
	)
	if err != nil {
		return nil, err
	}
	if len(usr_rows) == 0 {
		return nil, ErrGetUserById
	}
	usr := usr_rows[0]
	return &usr, nil
}

func (repo *PostgresRepository) GetAllUsers(ctx context.Context) ([]dto.UserDbo, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		SELECT id, state, email, name, lastname, avatar_url
		FROM users
		WHERE deleted_at IS NULL;
		`,
	)
	if err != nil {
		return nil, err
	}
	if len(usr_rows) == 0 {
		return nil, ErrGetUsers
	}
	return usr_rows, nil
}

func (repo *PostgresRepository) SetUserState(ctx context.Context, id int64, state models.State) error {
	_, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		UPDATE users
		SET state = $1::text
		WHERE id = $2::bigint
		  AND deleted_at IS NULL;
		`,
		string(state),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) SaveUser(ctx context.Context, user dto.UserDbo, user_id int64) (*dto.UserDbo, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		WITH data(id, email, name, lastname, avatar_url) AS (
		  SELECT $1::bigint, $2::text, $3::text, $4::text, $5::text
		)
		UPDATE users u
		SET email      = COALESCE(NULLIF(t.user_email, email), email),
			  name       = COALESCE(NULLIF(t.user_name, name), name),
			  lastname   = COALESCE(NULLIF(t.user_lastname, lastname), lastname),
			  avatar_url = COALESCE(NULLIF(t.user_avatar_url, avatar_url), avatar_url),
			  updated_at = NOW()
		FROM (
		SELECT d.id, d.email, d.name, d.lastname, d.avatar_url, $6::bigint, r.name
			FROM data d
			LEFT JOIN user_roles ur ON ur.user_id = $6::bigint
			JOIN roles r ON r.id = ur.role_id
		) AS t(user_id, user_email, user_name, user_lastname, user_avatar_url, update_user_id, role_name)
		WHERE u.id = t.user_id
			AND (u.id = t.update_user_id OR t.role_name = 'admin')
		RETURNING u.*;
		`,
		user.Id,
		user.Email,
		user.Name,
		user.LastName,
		user.AvatarUrl,
		user_id,
	)
	if err != nil {
		return nil, err
	}
	if len(usr_rows) == 0 {
		return nil, ErrUpdate
	}
	usr := usr_rows[0]
	return &usr, nil
}

func (repo *PostgresRepository) ValidateUserPermissions(ctx context.Context, user_id int64, perm models.Permission) error {
	perms, err := pu.Dispatch[dto.RolePermissionDbo](
		ctx,
		repo.db,
		`
		SELECT rp.*
		FROM role_permissions rp
		LEFT JOIN user_roles ur ON ur.role_id = rp.role_id
		LEFT JOIN users u ON u.id = ur.user_id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1::bigint
		  AND p.name = $2::text;
		`,
		user_id,
		perm.Name,
	)
	if err != nil {
		return err
	}
	if len(perms) == 0 {
		return ErrValidatePermissions
	}
	p := perms[0]
	if !p.Can(perm.Action) {
		return ErrValidatePermissions
	}
	return nil
}

func (repo *PostgresRepository) GrantRoleToUser(ctx context.Context, user_id int64, role_id int64) error {
	_, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		UPDATE user_roles
		SET role_id = $1::bigint
		WHERE user_id = $2::bigint
		`,
		role_id,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) DeleteUser(ctx context.Context, user_id int64) error {
	_, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		DELETE
		FROM users u
		WHERE u.id = $1::bigint
		  AND u.id IN (
      	SELECT ur.user_id
				FROM user_roles ur
				LEFT JOIN roles r ON r.id = ur.role_id
				WHERE ur.user_id = $1::bigint
				  AND r.name <> 'admin'
			);
		`,
		user_id,
	)
	if err != nil {
		return err
	}
	return nil
}
