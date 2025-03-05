package dto

import (
	"dashboard/internal/core/models"
)

type ChatMemberDto struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID                 int      `json:"id,omitempty"`
		FirstName          string   `json:"first_name,omitempty"`
		LastName           string   `json:"last_name,omitempty"`
		Username           string   `json:"username,omitempty"`
		Type               string   `json:"type,omitempty"`
		CanSendGift        bool     `json:"can_send_gift,omitempty"`
		ActiveUsernames    []string `json:"active_usernames,omitempty"`
		Bio                string   `json:"bio,omitempty"`
		HasPrivateForwards bool     `json:"has_private_forwards,omitempty"`
		BusinessIntro      struct {
			Title   string `json:"title,omitempty"`
			Message string `json:"message,omitempty"`
			Sticker struct {
				Width      int    `json:"width,omitempty"`
				Height     int    `json:"height,omitempty"`
				Emoji      string `json:"emoji,omitempty"`
				SetName    string `json:"set_name,omitempty"`
				IsAnimated bool   `json:"is_animated,omitempty"`
				IsVideo    bool   `json:"is_video,omitempty"`
				Type       string `json:"type,omitempty"`
				Thumbnail  struct {
					FileID       string `json:"file_id,omitempty"`
					FileUniqueID string `json:"file_unique_id,omitempty"`
					FileSize     int    `json:"file_size,omitempty"`
					Width        int    `json:"width,omitempty"`
					Height       int    `json:"height,omitempty"`
				} `json:"thumbnail,omitempty"`
				Thumb struct {
					FileID       string `json:"file_id,omitempty"`
					FileUniqueID string `json:"file_unique_id,omitempty"`
					FileSize     int    `json:"file_size,omitempty"`
					Width        int    `json:"width,omitempty"`
					Height       int    `json:"height,omitempty"`
				} `json:"thumb,omitempty"`
				FileID       string `json:"file_id,omitempty"`
				FileUniqueID string `json:"file_unique_id,omitempty"`
				FileSize     int    `json:"file_size,omitempty"`
			} `json:"sticker,omitempty"`
		} `json:"business_intro,omitempty"`
		BusinessLocation struct {
			Address string `json:"address,omitempty"`
		} `json:"business_location,omitempty"`
		BusinessOpeningHours struct {
			OpeningHours []struct {
				OpeningMinute int `json:"opening_minute,omitempty"`
				ClosingMinute int `json:"closing_minute,omitempty"`
			} `json:"opening_hours,omitempty"`
			TimeZoneName string `json:"time_zone_name,omitempty"`
		} `json:"business_opening_hours,omitempty"`
		Photo struct {
			SmallFileID       string `json:"small_file_id,omitempty"`
			SmallFileUniqueID string `json:"small_file_unique_id,omitempty"`
			BigFileID         string `json:"big_file_id,omitempty"`
			BigFileUniqueID   string `json:"big_file_unique_id,omitempty"`
		} `json:"photo,omitempty"`
		EmojiStatusCustomEmojiID string `json:"emoji_status_custom_emoji_id,omitempty"`
		MaxReactionCount         int    `json:"max_reaction_count,omitempty"`
		AccentColorID            int    `json:"accent_color_id,omitempty"`
		BackgroundCustomEmojiID  string `json:"background_custom_emoji_id,omitempty"`
	} `json:"result,omitempty"`
}

func (d ChatMemberDto) ToValue() models.ChatMember {
	return models.ChatMember{
		ID:        d.Result.ID,
		FirstName: d.Result.FirstName,
		LastName:  d.Result.LastName,
		Username:  d.Result.Username,
		Bio:       d.Result.Bio,
		Photo: struct {
			SmallFileID       string
			SmallFileUniqueID string
			BigFileID         string
			BigFileUniqueID   string
		}{
			SmallFileID:       d.Result.Photo.SmallFileID,
			SmallFileUniqueID: d.Result.Photo.SmallFileUniqueID,
			BigFileID:         d.Result.Photo.BigFileID,
			BigFileUniqueID:   d.Result.Photo.BigFileUniqueID,
		},
	}
}

type GetTelegramUserRequest struct {
	Id         int64 `json:"id"`
	BotId      int64 `json:"bot_id"`
	TelegramId int64 `json:"telegram_id"`
}

type SendTelegramMessageRequest struct {
	Id         int64  `json:"id"`
	BotId      int64  `json:"bot_id"`
	TelegramId int64  `json:"telegram_id"`
	Message    string `json:"message"`
}

type ChatFileDto struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		FileID       string `json:"file_id,omitempty"`
		FileUniqueID string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	} `json:"result,omitempty"`
}

type ChatSendMessageDto struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Username  string `json:"username,omitempty"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Username  string `json:"username,omitempty"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

type TelegramUserProfileDto struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Bio       string `json:"bio,omitempty"`
	PhotoUrl  string `json:"photo_url,omitempty"`
}

func ToTelegramUserProfileDto(m models.ChatMember) TelegramUserProfileDto {
	return TelegramUserProfileDto{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Username:  m.Username,
		Bio:       m.Bio,
		PhotoUrl:  m.PhotoUrl,
	}
}
