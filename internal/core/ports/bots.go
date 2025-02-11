package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type BotsService interface {
	GetBotById(context.Context, int64, int64) (models.Bot, error)
	GetAllBots(context.Context, int64) ([]models.Bot, error)
	CreateBot(context.Context, models.Bot, int64) (int64, error)
	UpdateBotById(context.Context, models.Bot, int64) (models.Bot, error)
	DeleteBotById(context.Context, int64, int64) error
	SetBotState(context.Context, models.BotState, int64, int64) error
	RunBot(context.Context, models.Bot, int64) error
	StopBot(context.Context, models.Bot, int64) error
}

type BotsRepository interface {
	GetBotById(context.Context, int64, int64) (*dto.BotDbo, error)
	GetAllBots(context.Context, int64) ([]dto.BotDbo, error)
	CreateBot(context.Context, dto.BotDbo, int64) (int64, error)
	UpdateBotById(context.Context, dto.BotDbo, int64) (*dto.BotDbo, error)
	DeleteBotById(context.Context, int64, int64) error
	SetBotState(context.Context, string, int64, int64) error
}

type BotsEventRepository interface {
	PublishBotsEvent(context.Context, models.Bot, int64, models.BotState) error
}
