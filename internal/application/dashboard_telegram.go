package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetTelegramUserById(user_id, statistics_id, bot_id, telegram_id int64) (models.ChatMember, error) {
	ctx := context.Background()
	d.log.Log("info", "getting telegram user by statistics_id", logger.WithInt64Attr("statistics_id", statistics_id))
	tg_user, err := d.telegram.GetTelegramUserById(
		ctx,
		user_id,
		statistics_id,
		bot_id,
		telegram_id,
	)
	if err != nil {
		d.log.Log("error", "error while getting telegram user by statistics_id", logger.WithErrAttr(err))
		return tg_user, err
	}
	d.log.Log("info", "telegram user successfully get", logger.WithInt64Attr("statistics_id", statistics_id))
	return tg_user, nil
}

func (d *DashboardService) SendMessage(user_id, statistics_id, bot_id, telegram_id int64, message string) error {
	ctx := context.Background()
	d.log.Log("info", "error while sending message to", logger.WithInt64Attr("telegram_id", telegram_id))
	if err := d.telegram.SendMessage(ctx, user_id, statistics_id, bot_id, telegram_id, message); err != nil {
		d.log.Log("error", "error while getting chat by chat_id", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "message successfully sent", logger.WithInt64Attr("telegram_id", telegram_id))
	return nil
}
