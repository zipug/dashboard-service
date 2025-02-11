package redis

import (
	"context"
	"dashboard/internal/core/models"
	"encoding/json"
	"errors"
	"time"
)

type BotsPayload struct {
	BotID       int64  `json:"bot_id"`
	ProjectID   int64  `json:"project_id"`
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (repo *RedisRepository) PublishBotsEvent(ctx context.Context, bot models.Bot, user_id int64, state models.BotState) error {
	var event_type string
	switch state {
	case models.RUNNING:
		event_type = "run"
	case models.STOPPED:
		event_type = "stop"
	default:
		return errors.New("unknown bot state")
	}
	evt := models.Event{
		Type:      event_type,
		Timestamp: time.Now().Unix(),
		Payload: BotsPayload{
			BotID:       bot.Id,
			ProjectID:   bot.ProjectId,
			UserID:      user_id,
			Name:        bot.Name,
			Description: bot.Description,
			Icon:        bot.Icon,
		},
	}
	message, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	if err := repo.rdb.Publish(ctx, "bot", message).Err(); err != nil {
		return err
	}
	return nil
}
