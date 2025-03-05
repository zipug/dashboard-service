package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetAllChats(user_id int64) ([]models.Chat, error) {
	ctx := context.Background()
	d.log.Log("info", "getting chats by user_id", logger.WithInt64Attr("user_id", user_id))
	chats, err := d.chat.GetAllChats(ctx, user_id)
	if err != nil {
		d.log.Log("error", "error while getting chats by user_id", logger.WithErrAttr(err))
		return chats, err
	}
	d.log.Log("info", "chats successfully get", logger.WithInt64Attr("user_id", user_id))
	return chats, nil
}

func (d *DashboardService) GetChatById(chat_id, user_id int64) ([]models.Chat, error) {
	ctx := context.Background()
	d.log.Log("info", "getting chat by chat_id", logger.WithInt64Attr("chat_id", chat_id))
	chats, err := d.chat.GetChatById(ctx, chat_id, user_id)
	if err != nil {
		d.log.Log("error", "error while getting chat by chat_id", logger.WithErrAttr(err))
		return chats, err
	}
	d.log.Log("info", "chat successfully get", logger.WithInt64Attr("chat_id", chat_id))
	return chats, nil
}

func (d *DashboardService) ResolveQuestion(chat_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "resolve chat by chat_id", logger.WithInt64Attr("chat_id", chat_id))
	if err := d.chat.ResolveQuestion(ctx, chat_id); err != nil {
		d.log.Log("error", "error while resolving chat by chat_id", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "chat successfully resolved", logger.WithInt64Attr("chat_id", chat_id))
	return nil
}
