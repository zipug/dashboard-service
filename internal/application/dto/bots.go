package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type BotDto struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	State       string `json:"state"`
}

type BotDbo struct {
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description,omitempty" json:"description,omitempty"`
	Icon        string         `db:"icon,omitempty" json:"icon,omitempty"`
	State       string         `db:"state" json:"state"`
	CreatedAt   sql.NullString `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   sql.NullString `db:"update_at,omitempty" json:"update_at,omitempty"`
	DeletedAt   sql.NullString `db:"delete_at,omitempty" json:"deleted_at,omitempty"`
}

func (dto *BotDto) ToValue() models.Bot {
	return models.Bot{
		Name:        dto.Name,
		Description: dto.Description,
		Icon:        dto.Icon,
		State:       dto.State,
	}
}

func (dbo *BotDbo) ToValue() models.Bot {
	return models.Bot{
		Name:        dbo.Name,
		Description: dbo.Description,
		Icon:        dbo.Icon,
		State:       dbo.State,
	}
}

func ToBotDto(m models.Bot) BotDto {
	return BotDto{
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		State:       m.State,
	}
}

func ToBotDbo(m models.Bot) BotDbo {
	return BotDbo{
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		State:       m.State,
	}
}
