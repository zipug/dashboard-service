package telegram

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TelegramService struct {
	repo TelegramRepo
}

type TelegramRepo interface {
	ports.TelegramRepository
	ports.BotsRepository
}

func NewTelegramService(repo TelegramRepo) *TelegramService {
	return &TelegramService{
		repo: repo,
	}
}

func (t *TelegramService) getChat(
	ctx context.Context,
	api_token string,
	telegram_id int64,
) (models.ChatMember, error) {
	var res models.ChatMember
	client := http.Client{}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	payload := strings.NewReader(fmt.Sprintf("chat_id=%d", telegram_id))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.telegram.org/bot%s/getChat", api_token), payload)
	if err != nil {
		return res, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	var response dto.ChatMemberDto
	if err := json.Unmarshal(readBytes, &response); err != nil {
		return res, err
	}
	if !response.Ok {
		return res, errors.New("error while getting chat info")
	}
	res = response.ToValue()
	return res, nil
}

func (t *TelegramService) getFile(
	ctx context.Context,
	api_token string,
	file_id string,
) (dto.ChatFileDto, error) {
	var res dto.ChatFileDto
	client := http.Client{}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	payload := strings.NewReader(fmt.Sprintf("file_id=%s", file_id))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.telegram.org/bot%s/getFile", api_token), payload)
	if err != nil {
		return res, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return res, err
	}
	if !res.Ok {
		return res, errors.New("error while getting file")
	}
	return res, nil
}

func (t *TelegramService) sendMessage(
	ctx context.Context,
	telegram_id int64,
	api_token string,
	message string,
) (dto.ChatSendMessageDto, error) {
	var res dto.ChatSendMessageDto
	client := http.Client{}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	payload := strings.NewReader(
		fmt.Sprintf(
			"entities=&reply_to_message_id=&disable_notification=&chat_id=%d&text=%s&parse_mode=markdown&allow_sending_without_reply=&disable_web_page_preview=&reply_markup=",
			telegram_id,
			message,
		),
	)
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", api_token), payload)
	if err != nil {
		return res, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return res, err
	}
	if !res.Ok {
		return res, errors.New("error while sending message")
	}
	return res, nil
}

func (t *TelegramService) GetTelegramUserById(
	ctx context.Context,
	user_id,
	statistics_id,
	bot_id,
	telegram_id int64,
) (
	models.ChatMember,
	error,
) {
	var res models.ChatMember
	botDbo, err := t.repo.GetBotById(ctx, bot_id, user_id)
	if err != nil {
		return res, err
	}
	if botDbo.ApiToken == "" {
		return res, errors.New("bot has not a valid APIToken")
	}
	chat, err := t.getChat(ctx, botDbo.ApiToken, telegram_id)
	if err != nil {
		return res, err
	}
	res = chat
	res.ApiToken = botDbo.ApiToken
	file, err := t.getFile(ctx, res.ApiToken, res.Photo.SmallFileID)
	if err != nil {
		return res, err
	}
	res.PhotoUrl = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", res.ApiToken, file.Result.FilePath)
	return res, nil
}

func (t *TelegramService) SendMessage(
	ctx context.Context,
	user_id,
	statistics_id,
	bot_id,
	telegram_id int64,
	message string,
) error {
	botDbo, err := t.repo.GetBotById(ctx, bot_id, user_id)
	if err != nil {
		return err
	}
	resp, err := t.sendMessage(ctx, telegram_id, botDbo.ApiToken, message)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return errors.New("error while sending telegram message")
	}
	return t.repo.CreateDialog(ctx, user_id, statistics_id, message)
}
