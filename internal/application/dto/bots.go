package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type BotDto struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	State       string `json:"state"`
}

type BotDbo struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description,omitempty" json:"description,omitempty"`
	Icon        string         `db:"icon,omitempty" json:"icon,omitempty"`
	State       string         `db:"state" json:"state"`
	UserID      int64          `db:"user_id" json:"user_id"`
	CreatedAt   sql.NullString `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   sql.NullString `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt   sql.NullString `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (dto *BotDto) ToValue() models.Bot {
	return models.Bot{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: dto.Description,
		Icon:        dto.Icon,
		State:       models.BotState(dto.State),
	}
}

func (dbo *BotDbo) ToValue() models.Bot {
	return models.Bot{
		Id:          dbo.Id,
		Name:        dbo.Name,
		Description: dbo.Description,
		Icon:        dbo.Icon,
		State:       models.BotState(dbo.State),
	}
}

func ToBotDto(m models.Bot) BotDto {
	return BotDto{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		State:       string(m.State),
	}
}

func ToBotDbo(m models.Bot) BotDbo {
	return BotDbo{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		State:       string(m.State),
	}
}
