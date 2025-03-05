package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type ChatsService interface {
	GetAllChats(ctx context.Context, user_id int64) ([]models.Chat, error)
	GetChatById(ctx context.Context, chat_id, user_id int64) ([]models.Chat, error)
	ResolveQuestion(ctx context.Context, chat_id int64) error
}

type ChatsRepository interface {
	GetAllChats(ctx context.Context, user_id int64) ([]dto.ChatDbo, error)
	GetChatById(ctx context.Context, chat_id, user_id int64) ([]dto.ChatDbo, error)
	ResolveQuestion(ctx context.Context, chat_id int64) error
}
