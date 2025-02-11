package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"errors"
)

func (d *DashboardService) GetBotById(bot_id, user_id int64) (models.Bot, error) {
	ctx := context.Background()
	d.log.Log("info", "get bot by id", logger.WithInt64Attr("bot_id", bot_id))
	bot, err := d.bot.GetBotById(ctx, bot_id, user_id)
	if err != nil {
		d.log.Log("error", "could not get bot by id", logger.WithInt64Attr("bot_id", bot_id))
		return bot, err
	}
	d.log.Log("info", "bot found by id", logger.WithInt64Attr("bot_id", bot_id))
	return bot, nil
}

func (d *DashboardService) GetAllBots(user_id int64) ([]models.Bot, error) {
	ctx := context.Background()
	d.log.Log("info", "get all bots by user_id", logger.WithInt64Attr("user_id", user_id))
	bots, err := d.bot.GetAllBots(ctx, user_id)
	if err != nil {
		d.log.Log("error", "could not get all bots by user_id", logger.WithInt64Attr("user_id", user_id))
		return nil, err
	}
	d.log.Log("info", "all bots found by user_id", logger.WithInt64Attr("user_id", user_id))
	return bots, nil
}

func (d *DashboardService) CreateBot(bot models.Bot, user_id int64) (int64, error) {
	ctx := context.Background()
	d.log.Log("info", "create bot", logger.WithInt64Attr("user_id", user_id))
	id, err := d.bot.CreateBot(ctx, bot, user_id)
	if err != nil {
		d.log.Log("error", "could not create bot", logger.WithInt64Attr("user_id", user_id))
		return id, err
	}
	d.log.Log("info", "bot created", logger.WithInt64Attr("user_id", user_id))
	return id, nil
}

func (d *DashboardService) UpdateBotById(bot models.Bot, user_id int64) (models.Bot, error) {
	ctx := context.Background()
	d.log.Log("info", "update bot by id", logger.WithInt64Attr("bot_id", bot.Id))
	new_bot, err := d.bot.UpdateBotById(ctx, bot, user_id)
	if err != nil {
		return new_bot, err
	}
	d.log.Log("info", "bot updated by id", logger.WithInt64Attr("bot_id", bot.Id))
	return new_bot, nil
}

func (d *DashboardService) DeleteBotById(bot_id, user_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "delete bot by id", logger.WithInt64Attr("bot_id", bot_id))
	if err := d.bot.DeleteBotById(ctx, bot_id, user_id); err != nil {
		d.log.Log("error", "could not delete bot by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	if err := d.bot.SetBotState(ctx, models.DELETED, bot_id, user_id); err != nil {
		d.log.Log("error", "could not set bot state by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	d.log.Log("info", "bot deleted by id", logger.WithInt64Attr("bot_id", bot_id))
	return nil
}

func (d *DashboardService) RunBotById(bot_id, user_id int64) error {
	ctx := context.Background()
	bot, err := d.GetBotById(bot_id, user_id)
	if err != nil {
		return err
	}
	if bot.State == models.RUNNING || bot.State == models.DELETED {
		d.log.Log("error", "bot is already running or deleted", logger.WithInt64Attr("bot_id", bot_id))
		return errors.New("bot is already running or deleted")
	}
	d.log.Log("info", "set bot state by id", logger.WithInt64Attr("bot_id", bot_id))
	if err := d.bot.SetBotState(ctx, models.RUNNING, bot_id, user_id); err != nil {
		d.log.Log("error", "could not set bot state by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	d.log.Log("info", "bot state set by id", logger.WithInt64Attr("bot_id", bot_id))
	d.log.Log("info", "run bot by id", logger.WithInt64Attr("bot_id", bot_id))
	if err := d.bot.RunBot(ctx, bot, user_id); err != nil {
		d.log.Log("error", "could not run bot by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	d.log.Log("info", "bot successfully runned by id", logger.WithInt64Attr("bot_id", bot_id))
	return nil
}

func (d *DashboardService) StopBotById(bot_id, user_id int64) error {
	ctx := context.Background()
	bot, err := d.GetBotById(bot_id, user_id)
	if err != nil {
		return err
	}
	if bot.State == models.STOPPED || bot.State == models.DELETED {
		d.log.Log("error", "bot is already stopped or deleted", logger.WithInt64Attr("bot_id", bot_id))
		return errors.New("bot is already stopped or deleted")
	}
	d.log.Log("info", "set bot state by id", logger.WithInt64Attr("bot_id", bot_id))
	if err := d.bot.SetBotState(ctx, models.STOPPED, bot_id, user_id); err != nil {
		d.log.Log("error", "could not set bot state by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	d.log.Log("info", "bot state set by id", logger.WithInt64Attr("bot_id", bot_id))
	d.log.Log("info", "stop bot by id", logger.WithInt64Attr("bot_id", bot_id))
	if err := d.bot.StopBot(ctx, bot, user_id); err != nil {
		d.log.Log("error", "could not stop bot by id", logger.WithInt64Attr("bot_id", bot_id))
		return err
	}
	d.log.Log("info", "bot successfully stopped by id", logger.WithInt64Attr("bot_id", bot_id))
	return nil
}
