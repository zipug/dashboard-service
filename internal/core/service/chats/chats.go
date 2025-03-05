package chats

import (
	"context"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type ChatsService struct {
	repo ports.ChatsRepository
}

func NewChatsService(repo ports.ChatsRepository) *ChatsService {
	return &ChatsService{repo: repo}
}

func (c *ChatsService) GetChatById(ctx context.Context, chat_id, user_id int64) ([]models.Chat, error) {
	dbo, err := c.repo.GetChatById(ctx, chat_id, user_id)
	if err != nil {
		return nil, err
	}
	var chats []models.Chat
	for _, chat := range dbo {
		chats = append(chats, chat.ToValue())
	}
	return chats, nil
}

func (c *ChatsService) GetAllChats(ctx context.Context, user_id int64) ([]models.Chat, error) {
	dbo, err := c.repo.GetAllChats(ctx, user_id)
	if err != nil {
		return []models.Chat{}, err
	}
	var chats []models.Chat
	for _, chat := range dbo {
		chats = append(chats, chat.ToValue())
	}
	return chats, nil
}

func (c *ChatsService) ResolveQuestion(ctx context.Context, chat_id int64) error {
	return c.repo.ResolveQuestion(ctx, chat_id)
}
