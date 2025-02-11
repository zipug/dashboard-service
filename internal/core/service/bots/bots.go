package bots

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type BotsService struct {
	repo  ports.BotsRepository
	event ports.BotsEventRepository
}

func NewBotsService(repo ports.BotsRepository, event ports.BotsEventRepository) *BotsService {
	return &BotsService{repo: repo, event: event}
}

func (s *BotsService) GetBotById(ctx context.Context, bot_id, user_id int64) (models.Bot, error) {
	dbo, err := s.repo.GetBotById(ctx, bot_id, user_id)
	if err != nil {
		return models.Bot{}, err
	}
	return dbo.ToValue(), nil
}

func (s *BotsService) GetAllBots(ctx context.Context, user_id int64) ([]models.Bot, error) {
	dbo, err := s.repo.GetAllBots(ctx, user_id)
	if err != nil {
		return nil, err
	}
	var res []models.Bot
	for _, d := range dbo {
		res = append(res, d.ToValue())
	}
	return res, nil
}

func (s *BotsService) CreateBot(ctx context.Context, bot models.Bot, user_id int64) (int64, error) {
	bot.State = models.CREATED
	dbo := dto.ToBotDbo(bot)
	id, err := s.repo.CreateBot(ctx, dbo, user_id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (s *BotsService) UpdateBotById(ctx context.Context, bot models.Bot, user_id int64) (models.Bot, error) {
	dbo := dto.ToBotDbo(bot)
	new_dbo, err := s.repo.UpdateBotById(ctx, dbo, user_id)
	if err != nil {
		return models.Bot{}, err
	}
	return new_dbo.ToValue(), nil
}

func (s *BotsService) DeleteBotById(ctx context.Context, bot_id, user_id int64) error {
	return s.repo.DeleteBotById(ctx, bot_id, user_id)
}

func (s *BotsService) SetBotState(ctx context.Context, state models.BotState, bot_id, user_id int64) error {
	return s.repo.SetBotState(ctx, string(state), bot_id, user_id)
}

func (s *BotsService) RunBot(ctx context.Context, bot models.Bot, user_id int64) error {
	return s.event.PublishBotsEvent(ctx, bot, user_id, models.RUNNING)
}

func (s *BotsService) StopBot(ctx context.Context, bot models.Bot, user_id int64) error {
	return s.event.PublishBotsEvent(ctx, bot, user_id, models.STOPPED)
}
