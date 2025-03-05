package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
	"time"
)

type ChatDbo struct {
	Id         int64         `db:"id"`
	BotId      int64         `db:"bot_id"`
	TelegramId int64         `db:"telegram_id"`
	ProjectId  int64         `db:"project_id"`
	CreatedBy  int64         `db:"created_by"`
	UserId     int64         `db:"user_id"`
	Name       string        `db:"name"`
	Question   string        `db:"question"`
	CreatedAt  sql.NullTime  `db:"created_at"`
	IsResolved bool          `db:"is_resolved"`
	ParentId   sql.NullInt64 `db:"parent_id"`
}

type StatisticDbo struct {
	Id          int64          `db:"id"`
	BotId       int64          `db:"bot_id"`
	TelegramId  int64          `db:"telegram_id"`
	Question    sql.NullString `db:"question"`
	ArticleId   sql.NullInt64  `db:"article_id"`
	ArticleName string         `db:"article_name"`
	IsResolved  bool           `db:"is_resolved"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	ParentId    sql.NullInt64  `db:"parent_id"`
}

type ChatDto struct {
	Id         int64     `json:"id"`
	BotId      int64     `json:"bot_id"`
	TelegramId int64     `json:"telegram_id"`
	ProjectId  int64     `json:"project_id"`
	CreatedBy  int64     `json:"created_by,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	Name       string    `json:"name"`
	Question   string    `json:"question"`
	CreatedAt  time.Time `json:"created_at"`
	IsResolved bool      `json:"is_resolved"`
	ParentId   int64     `json:"parent_id"`
}

func (c ChatDbo) ToValue() models.Chat {
	return models.Chat{
		Id:         c.Id,
		BotId:      c.BotId,
		TelegramId: c.TelegramId,
		ProjectId:  c.ProjectId,
		CreatedBy:  c.CreatedBy,
		UserId:     c.UserId,
		Name:       c.Name,
		Question:   c.Question,
		CreatedAt:  c.CreatedAt.Time,
		IsResolved: c.IsResolved,
		ParentId:   c.ParentId.Int64,
	}
}

func (c ChatDto) ToValue() models.Chat {
	return models.Chat{
		Id:         c.Id,
		BotId:      c.BotId,
		TelegramId: c.TelegramId,
		ProjectId:  c.ProjectId,
		CreatedBy:  c.CreatedBy,
		UserId:     c.UserId,
		Name:       c.Name,
		Question:   c.Question,
		CreatedAt:  c.CreatedAt,
		IsResolved: c.IsResolved,
		ParentId:   c.ParentId,
	}
}

func ToChatDbo(c models.Chat) ChatDbo {
	return ChatDbo{
		Id:         c.Id,
		BotId:      c.BotId,
		TelegramId: c.TelegramId,
		ProjectId:  c.ProjectId,
		CreatedBy:  c.CreatedBy,
		UserId:     c.UserId,
		Name:       c.Name,
		Question:   c.Question,
		CreatedAt:  sql.NullTime{Time: c.CreatedAt, Valid: true},
		IsResolved: c.IsResolved,
		ParentId:   sql.NullInt64{Int64: c.ParentId, Valid: true},
	}
}

func ToChatDto(c models.Chat) ChatDto {
	return ChatDto{
		Id:         c.Id,
		BotId:      c.BotId,
		TelegramId: c.TelegramId,
		ProjectId:  c.ProjectId,
		CreatedBy:  c.CreatedBy,
		UserId:     c.UserId,
		Name:       c.Name,
		Question:   c.Question,
		CreatedAt:  c.CreatedAt,
		IsResolved: c.IsResolved,
		ParentId:   c.ParentId,
	}
}
