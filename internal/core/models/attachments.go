package models

type Attachment struct {
	Id          int64
	Name        string
	Description string
	UserId      int64
	URL         string
}
