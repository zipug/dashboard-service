package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type BotDto struct {
	Id          int64  `json:"id"`
	ProjectId   int64  `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	State       string `json:"state"`
	ApiToken    string `json:"api_token"`
}

type BotDbo struct {
	Id          int64          `db:"id" json:"id"`
	ProjectId   int64          `db:"project_id" json:"project_id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description,omitempty" json:"description,omitempty"`
	Icon        sql.NullString `db:"icon,omitempty" json:"icon,omitempty"`
	State       string         `db:"state" json:"state"`
	UserID      int64          `db:"user_id" json:"user_id"`
	ApiToken    string         `db:"api_token" json:"api_token"`
	CreatedAt   sql.NullString `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   sql.NullString `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt   sql.NullString `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (dto *BotDto) ToValue() models.Bot {
	return models.Bot{
		Id:          dto.Id,
		ProjectId:   dto.ProjectId,
		Name:        dto.Name,
		Description: dto.Description,
		Icon:        dto.Icon,
		State:       models.BotState(dto.State),
		ApiToken:    dto.ApiToken,
	}
}

func (dbo *BotDbo) ToValue() models.Bot {
	return models.Bot{
		Id:          dbo.Id,
		ProjectId:   dbo.ProjectId,
		Name:        dbo.Name,
		Description: dbo.Description.String,
		Icon:        dbo.Icon.String,
		State:       models.BotState(dbo.State),
		ApiToken:    dbo.ApiToken,
	}
}

func ToBotDto(m models.Bot) BotDto {
	return BotDto{
		Id:          m.Id,
		ProjectId:   m.ProjectId,
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		State:       string(m.State),
		ApiToken:    m.ApiToken,
	}
}

func ToBotDbo(m models.Bot) BotDbo {
	return BotDbo{
		Id:          m.Id,
		ProjectId:   m.ProjectId,
		Name:        m.Name,
		Description: sql.NullString{String: m.Description, Valid: true},
		Icon:        sql.NullString{String: m.Icon, Valid: true},
		State:       string(m.State),
		ApiToken:    m.ApiToken,
	}
}
