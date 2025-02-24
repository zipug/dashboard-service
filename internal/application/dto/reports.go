package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type ReportDto struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

type ReportDbo struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Description string         `db:"description,omitempty" json:"description,omitempty"`
	Icon        string         `db:"icon,omitempty" json:"icon,omitempty"`
	CreatedAt   sql.NullString `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   sql.NullString `db:"update_at,omitempty" json:"update_at,omitempty"`
	DeletedAt   sql.NullString `db:"delete_at,omitempty" json:"deleted_at,omitempty"`
}

func (dto *ReportDto) ToValue() models.Report {
	return models.Report{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: dto.Description,
		Icon:        dto.Icon,
	}
}

func (dbo *ReportDbo) ToValue() models.Report {
	return models.Report{
		Id:          dbo.Id,
		Name:        dbo.Name,
		Description: dbo.Description,
		Icon:        dbo.Icon,
	}
}

func ToReportDto(m models.Report) ReportDto {
	return ReportDto{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
	}
}

func ToReportDbo(m models.Report) ReportDbo {
	return ReportDbo{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
	}
}
