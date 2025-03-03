package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

var BadUserId = -1

type AuthenticateDto struct {
	Token string `json:"token,omitempty"`
}

type UserDto struct {
	Id             int64        `json:"id,omitempty"`
	State          models.State `json:"state,omitempty"`
	Email          string       `json:"email,omitempty"`
	Password       string       `json:"password,omitempty"`
	RepeatPassword string       `json:"repeat_password,omitempty"`
	Name           string       `json:"name,omitempty"`
	LastName       string       `json:"lastname,omitempty"`
	AvatarUrl      string       `json:"avatar_url,omitempty"`
	Role           RoleDto      `json:"role,omitempty"`
	CreatedBy      int64        `json:"created_by,omitempty"`
}

//go:generate ../../../cmd/generator/main.go
type UpdateUserDto struct {
	Id        int64  `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

type UserDbo struct {
	Id        int64          `db:"id"`
	State     models.State   `db:"state"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Name      string         `db:"name"`
	LastName  string         `db:"lastname"`
	AvatarUrl sql.NullString `db:"avatar_url"`
	CreatedAt sql.NullTime   `db:"created_at,omitempty"`
	UpdatedAt sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt  sql.NullTime   `db:"deleted_at,omitempty"`
	CreatedBy sql.NullInt64  `db:"created_by,omitempty"`
}

type SafeUserDto struct {
	State     models.State `json:"state,omitempty"`
	Email     string       `json:"email,omitempty"`
	Name      string       `json:"name,omitempty"`
	LastName  string       `json:"lastname,omitempty"`
	AvatarUrl string       `json:"avatar_url,omitempty"`
	Role      RoleDto      `json:"role,omitempty"`
	CreatedBy int64        `json:"created_by,omitempty"`
}

type VerifyUserDto struct {
	Code models.OTPCode `json:"code"`
}

func (dto *UserDto) ToValue() models.User {
	return models.User{
		Id:             models.Id(dto.Id),
		State:          models.State(dto.State),
		Email:          models.Email(dto.Email),
		Password:       models.Password(dto.Password),
		RepeatPassword: models.Password(dto.RepeatPassword),
		Name:           models.Name(dto.Name),
		LastName:       models.LastName(dto.LastName),
		AvatarUrl:      models.AvatarUrl(dto.AvatarUrl),
		CreatedBy:      models.CreatedBy(dto.CreatedBy),
	}
}

func (dto *UpdateUserDto) ToValue() models.User {
	return models.User{
		Id:        models.Id(dto.Id),
		Email:     models.Email(dto.Email),
		Name:      models.Name(dto.Name),
		LastName:  models.LastName(dto.LastName),
		AvatarUrl: models.AvatarUrl(dto.AvatarUrl),
	}
}

func (dbo *UserDbo) ToValue() models.User {
	return models.User{
		Id:        models.Id(dbo.Id),
		State:     models.State(dbo.State),
		Email:     models.Email(dbo.Email),
		Password:  models.Password(dbo.Password),
		Name:      models.Name(dbo.Name),
		LastName:  models.LastName(dbo.LastName),
		AvatarUrl: models.AvatarUrl(dbo.AvatarUrl.String),
		CreatedBy: models.CreatedBy(dbo.CreatedBy.Int64),
	}
}

func ToUserDto(user models.User) UserDto {
	return UserDto{
		Id:             int64(user.Id),
		State:          user.State,
		Email:          string(user.Email),
		Password:       string(user.Password),
		RepeatPassword: string(user.RepeatPassword),
		Name:           string(user.Name),
		LastName:       string(user.LastName),
		AvatarUrl:      string(user.AvatarUrl),
		CreatedBy:      int64(user.CreatedBy),
	}
}

func ToSafeUserDto(user models.User) SafeUserDto {
	return SafeUserDto{
		State:     user.State,
		Email:     string(user.Email),
		Name:      string(user.Name),
		LastName:  string(user.LastName),
		AvatarUrl: string(user.AvatarUrl),
		CreatedBy: int64(user.CreatedBy),
	}
}

func ToUserDbo(user models.User) UserDbo {
	return UserDbo{
		Id:        int64(user.Id),
		State:     user.State,
		Email:     string(user.Email),
		Password:  string(user.Password),
		Name:      string(user.Name),
		LastName:  string(user.LastName),
		AvatarUrl: sql.NullString{String: string(user.AvatarUrl), Valid: true},
		CreatedBy: sql.NullInt64{Int64: int64(user.CreatedBy), Valid: true},
	}
}
