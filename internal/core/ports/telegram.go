package ports

import (
	"context"
	"dashboard/internal/core/models"
)

type TelegramService interface {
	GetTelegramUserById(ctx context.Context, user_id, statistics_id, bot_id, telegram_id int64) (models.ChatMember, error)
	SendMessage(ctx context.Context, user_id, statistics_id, bot_id, telegram_id int64, message string) error
}

type TelegramRepository interface {
	CreateDialog(ctx context.Context, user_id, statistics_id int64, answer string) error
}
