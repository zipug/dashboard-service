package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type UserService interface {
	RegisterUser(context.Context, models.User) (int64, error)
	LoginUser(context.Context, models.User) (models.User, error)
	GetUserById(context.Context, int64) (models.User, error)
	GetAllUsers(context.Context) ([]models.User, error)
	SaveUser(context.Context, models.User) error
	DeleteUser(context.Context, int64) error
}

type UserRepository interface {
	RegisterUser(context.Context, dto.UserDbo) (int64, error)
	GetUserById(context.Context, int64) (*dto.UserDbo, error)
	GetUserByEmail(context.Context, string) (*dto.UserDbo, error)
	GetAllUsers(context.Context) ([]dto.UserDbo, error)
	SaveUser(context.Context, dto.UserDbo) error
	DeleteUser(context.Context, int64) error
}

type UserDto interface {
	ToValue() models.User
}

type UserDbo interface {
	Tovalue() models.User
}
