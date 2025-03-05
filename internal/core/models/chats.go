package models

import "time"

type Chat struct {
	Id         int64
	BotId      int64
	TelegramId int64
	ProjectId  int64
	CreatedBy  int64
	UserId     int64
	Name       string
	Question   string
	ParentId   int64
	CreatedAt  time.Time
	IsResolved bool
}
