package redis

import (
	"context"
	"dashboard/internal/core/models"
	"encoding/json"
	"time"
)

type ReportsPayload struct {
	ReportID int64  `json:"report_id"`
	UserID   int64  `json:"user_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (repo *RedisRepository) ExecuteReport(ctx context.Context, report_id, user_id int64, date_from, date_to string) error {
	evt := models.Event{
		Type:      "report",
		Timestamp: time.Now().Unix(),
		Payload: ReportsPayload{
			ReportID: report_id,
			UserID:   report_id,
			DateFrom: date_from,
			DateTo:   date_to,
		},
	}
	message, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	if err := repo.rdb.Publish(ctx, "report", message).Err(); err != nil {
		return err
	}
	return nil
}
