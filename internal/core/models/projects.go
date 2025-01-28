package models

type Project struct {
	Id          int64
	Name        string
	Description string
	AvatarUrl   string
	RemoteUrl   string
	UserId      int64
}

type ProjectsContent struct {
	Project Project
	Content []Article
}
