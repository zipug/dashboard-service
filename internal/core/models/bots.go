package models

type Bot struct {
	Id          int64
	Name        string
	Description string
	Icon        string
	State       BotState
}

type BotState string

const (
	CREATED BotState = "created"
	RUNNING BotState = "running"
	STOPPED BotState = "stopped"
	DELETED BotState = "deleted"
)
