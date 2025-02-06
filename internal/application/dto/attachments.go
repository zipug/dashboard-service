package dto

import (
	"dashboard/internal/core/models"
	"database/sql"
)

type AttachmentDto struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Mimetype    string `json:"mimetype,omitempty"`
}

type AttachmentArticleDbo struct {
	AttachmentId int64          `db:"attachment_id"`
	ArticleId    int64          `db:"article_id"`
	CreatedAt    sql.NullString `db:"created_at,omitempty"`
	UpdateAt     sql.NullString `db:"updated_at,omitempty"`
	DeleteAt     sql.NullString `db:"deleted_at,omitempty"`
}

type AttachmentDbo struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	ObjectId    string         `db:"object_id"`
	Mimetype    string         `db:"mimetype"`
	UserID      int64          `db:"user_id,omitempty"`
	CreatedAt   sql.NullTime   `db:"created_at,omitempty"`
	UpdateAt    sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt    sql.NullTime   `db:"deleted_at,omitempty"`
}

type AttachmentNullSafeDbo struct {
	Id          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ObjectId    string `db:"object_id" json:"object_id"`
	Mimetype    string `db:"mimetype" json:"mimetype"`
	UserID      int64  `db:"user_id,omitempty" json:"user_id,omitempty"`
	CreatedAt   string `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdateAt    string `db:"updated_at,omitempty" json:"update_at,omitempty"`
	DeleteAt    string `db:"deleted_at,omitempty" json:"delete_at,omitempty"`
}

func (a *AttachmentNullSafeDbo) ToDbo() AttachmentDbo {
	return AttachmentDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		ObjectId:    a.ObjectId,
		Mimetype:    a.Mimetype,
	}
}

func (a *AttachmentNullSafeDbo) ToDto() AttachmentDto {
	return AttachmentDto{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description,
	}
}

func (a *AttachmentDbo) ToDto() AttachmentDto {
	return AttachmentDto{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		Mimetype:    a.Mimetype,
	}
}

func (a *AttachmentDbo) ToValue() models.Attachment {
	return models.Attachment{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		UserId:      a.UserID,
		ObjectId:    a.ObjectId,
		Mimetype:    a.Mimetype,
	}
}

func (a *AttachmentDbo) ToNullSafeDbo() AttachmentNullSafeDbo {
	return AttachmentNullSafeDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		ObjectId:    a.ObjectId,
		Mimetype:    a.Mimetype,
	}
}

func (a *AttachmentDto) ToDbo() AttachmentDbo {
	return AttachmentDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
	}
}

func ToAttachmentDto(a models.Attachment) AttachmentDto {
	return AttachmentDto{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		Mimetype:    a.Mimetype,
		URL:         a.URL,
	}
}

func ToAttachmentDbo(a models.Attachment) AttachmentDbo {
	return AttachmentDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		UserID:      a.UserId,
		ObjectId:    a.ObjectId,
		Mimetype:    a.Mimetype,
	}
}
