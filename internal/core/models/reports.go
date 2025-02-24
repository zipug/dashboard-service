package models

type Report struct {
	Id          int64
	Name        string
	Description string
	Icon        string
}

const REPORT_EVENT = "report_event"
