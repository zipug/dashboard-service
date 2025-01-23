package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/models"
	pu "dashboard/pkg/postgres_utils"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrConnect              = errors.New("could not connect to database")
	ErrPing                 = errors.New("could not ping database")
	ErrRegister             = errors.New("could not register user")
	ErrUpdate               = errors.New("could not update user")
	ErrPqEmailAlreadyExists = errors.New("pq: duplicate key value violates unique constraint \"users_email_unique\"")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrGetUserById          = errors.New("could not get user by id")
	ErrGetUsers             = errors.New("could not get users")
)

type PostgresRepository struct {
	db             *sqlx.DB
	uri            string
	migrationsPath string
	version        uint
	dirty          bool
}

func NewPostgresRepository(cfg *config.AppConfig) *PostgresRepository {
	repo := &PostgresRepository{}
	if err := repo.InvokeConnect(cfg); err != nil {
		panic(err)
	}
	return repo
}

func (repo *PostgresRepository) InvokeConnect(cfg *config.AppConfig) error {
	postgres_uri := pu.ParseURI(
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		cfg.Postgres.Port,
	)
	repo.uri = postgres_uri
	repo.migrationsPath = cfg.Postgres.MigrationsPath

	db, err := sqlx.Open("postgres", postgres_uri)
	if err != nil {
		return ErrConnect
	}
	repo.db = db
	if err := repo.PingTest(); err != nil {
		panic(err)
	}
	return nil
}

func (repo *PostgresRepository) PingTest() error {
	max_errs := 5
	errs := 0
	timeout := 1 * time.Second
	for max_errs > 0 {
		if err := repo.db.Ping(); err != nil {
			fmt.Printf("could not ping database: %s\n", err.Error())
			fmt.Printf("retrying in %s\n", timeout)
			max_errs--
			errs++
			time.Sleep(timeout)
		}
		max_errs = 0
		errs = 0
	}
	if errs == 0 {
		return nil
	}
	return fmt.Errorf("%w: postgres_uri: %s", ErrPing, repo.uri)
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) Migrate() error {
	if status := pu.Migrate(repo.uri, repo.migrationsPath, pu.Up); status.Error != nil {
		if status.Error != pu.ErrNoChange {
			return status.Error
		}
	} else {
		repo.version = status.Version
		repo.dirty = status.Dirty
	}
	return nil
}

func (repo *PostgresRepository) RegisterUser(ctx context.Context, user dto.UserDbo) (int64, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
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
		if errors.Is(err, ErrPqEmailAlreadyExists) {
			return int64(dto.BadUserId), ErrEmailAlreadyExists
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_unique\"" {
			return int64(dto.BadUserId), ErrEmailAlreadyExists
		}
		return int64(dto.BadUserId), err
	}
	if len(usr_rows) == 0 {
		return int64(dto.BadUserId), ErrRegister
	}
	usr := usr_rows[0]
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

func (repo *PostgresRepository) SaveUser(ctx context.Context, user dto.UserDbo) (*dto.UserDbo, error) {
	usr_rows, err := pu.Dispatch[dto.UserDbo](
		ctx,
		repo.db,
		`
		UPDATE users
		SET email = $1::text,
		    name = $2::text,
		    lastname = $3::text,
		    avatar_url = $4::text
		WHERE id = $5::bigint
		  AND deleted_at IS NULL
		RETURNING *;
		`,
		user.Email,
		user.Name,
		user.LastName,
		user.AvatarUrl,
		user.Id,
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

func (repo *PostgresRepository) DeleteUser(ctx context.Context, id int64) error {
	// delete user
	if id == -1 {
		return fmt.Errorf("id is required")
	}
	return nil
}
