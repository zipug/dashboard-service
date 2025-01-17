package ports

import "dashboard/internal/core/models"

type OTPService interface {
	SendOTPByUserId(int64, models.OTPTarget) error
	SendOTPByUserEmail(string, models.OTPTarget) error
}
