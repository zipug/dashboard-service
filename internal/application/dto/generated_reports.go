package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
	"encoding/json"
	"slices"
)

type GeneratedReportDto struct {
	Id       int64             `json:"id"`
	UserID   int64             `json:"user_id"`
	ReportID int64             `json:"report_id"`
	DateFrom string            `json:"date_from"`
	DateTo   string            `json:"date_to"`
	ObjectID string            `json:"object_id"`
	Content  map[string]string `json:"content"`
}

type GeneratedReportNullContentDto struct {
	UserID   int64  `json:"user_id"`
	ReportID int64  `json:"report_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	ObjectID string `json:"object_id"`
}

type GeneratedReportDbo struct {
	Id        int64          `db:"id" json:"id,omitempty"`
	UserID    int64          `db:"user_id" json:"user_id,omitempty"`
	ReportID  int64          `db:"report_id" json:"report_id,omitempty"`
	DateFrom  string         `db:"date_from" json:"date_from,omitempty"`
	DateTo    string         `db:"date_to" json:"date_to,omitempty"`
	ObjectID  string         `db:"object_id" json:"object_id,omitempty"`
	Content   sql.NullString `db:"content" json:"content,omitempty"`
	CreatedAt sql.NullTime   `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt sql.NullTime   `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (d *GeneratedReportDbo) ToValue() (models.GeneratedReport, error) {
	var content []map[string]string
	var headers []string
	if err := json.Unmarshal([]byte(d.Content.String), &content); err != nil {
		return models.GeneratedReport{}, err
	}
	rows := make([][]string, len(content))
	for i, item := range content {
		for key, value := range item {
			if !slices.Contains(headers, key) {
				headers = append(headers, key)
			}
			rows[i] = append(rows[i], value)
		}
	}
	return models.GeneratedReport{
		Id:       d.Id,
		UserID:   d.UserID,
		ReportID: d.ReportID,
		DateFrom: d.DateFrom,
		DateTo:   d.DateTo,
		ObjectID: d.ObjectID,
		Content: models.GeneratedReportContent{
			Headers: headers,
			Rows:    rows,
		},
	}, nil
}

func (d *GeneratedReportDbo) ToNullContentValue() models.GeneratedReport {
	return models.GeneratedReport{
		Id:       d.Id,
		UserID:   d.UserID,
		ReportID: d.ReportID,
		DateFrom: d.DateFrom,
		DateTo:   d.DateTo,
		ObjectID: d.ObjectID,
	}
}

func (d *GeneratedReportDto) ToValue() models.GeneratedReport {
	return models.GeneratedReport{
		Id:       d.Id,
		UserID:   d.UserID,
		ReportID: d.ReportID,
		DateFrom: d.DateFrom,
		DateTo:   d.DateTo,
		ObjectID: d.ObjectID,
	}
}

func ToGeneratedReportDto(value models.GeneratedReport) GeneratedReportDto {
	return GeneratedReportDto{
		Id:       value.Id,
		UserID:   value.UserID,
		ReportID: value.ReportID,
		DateFrom: value.DateFrom,
		DateTo:   value.DateTo,
		ObjectID: value.ObjectID,
	}
}

func ToGeneratedReportNullContentDto(value models.GeneratedReport) GeneratedReportNullContentDto {
	return GeneratedReportNullContentDto{
		UserID:   value.UserID,
		ReportID: value.ReportID,
		DateFrom: value.DateFrom,
		DateTo:   value.DateTo,
		ObjectID: value.ObjectID,
	}
}

func ToGeneratedReportDbo(value models.GeneratedReport) GeneratedReportDbo {
	return GeneratedReportDbo{
		Id:       value.Id,
		UserID:   value.UserID,
		ReportID: value.ReportID,
		DateFrom: value.DateFrom,
		DateTo:   value.DateTo,
		ObjectID: value.ObjectID,
	}
}
