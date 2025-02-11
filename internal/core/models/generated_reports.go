package models

type GeneratedReport struct {
	Id       int64
	UserID   int64
	ReportID int64
	DateFrom string
	DateTo   string
	ObjectID string
	Content  GeneratedReportContent
}

type GeneratedReportContent struct {
	Headers []string
	Rows    [][]string
}
