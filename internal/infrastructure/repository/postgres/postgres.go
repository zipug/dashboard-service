package postgres

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/common/service/config"
	pu "dashboard/pkg/postgres_utils"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrConnect              = errors.New("could not connect to database")
	ErrPing                 = errors.New("could not ping database")
	ErrRegister             = errors.New("could not register user")
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
	repo.InvokeConnect(cfg)
	fmt.Printf("Postgres connected on port: %v\n", cfg)
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
	if err := db.Ping(); err != nil {
		return fmt.Errorf("%w: postgres_uri: %s", ErrPing, postgres_uri)
	}
	repo.db = db
	return nil
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

func (repo *PostgresRepository) SaveUser(ctx context.Context, user dto.UserDbo) error {
	return nil
}

func (repo *PostgresRepository) DeleteUser(ctx context.Context, id int64) error {
	// delete user
	if id == -1 {
		return fmt.Errorf("id is required")
	}
	return nil
}
