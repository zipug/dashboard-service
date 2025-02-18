package models

type Bot struct {
	Id          int64
	ProjectId   int64
	Name        string
	Description string
	Icon        string
	State       BotState
	ApiToken    string
}

type BotState string

const (
	CREATED BotState = "created"
	RUNNING BotState = "running"
	STOPPED BotState = "stopped"
	DELETED BotState = "deleted"
)
