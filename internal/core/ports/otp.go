package ports

import (
	"context"
	"dashboard/internal/core/models"
)

type OTPService interface {
	SendOTP(context.Context, int64) error
	GetOTP(context.Context, int64) (models.OTPCode, error)
}

type OTPRepository interface {
	SaveUserOTP(context.Context, int64, models.OTPCode) error
	GetUserOTP(context.Context, int64) (models.OTPCode, error)
}
