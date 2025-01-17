package otp

import (
	"crypto/rand"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/models"
	"io"
	"strconv"
)

type OTPService struct {
	Max int
}

func NewOTPService(cfg *config.AppConfig) *OTPService {
	return &OTPService{Max: cfg.OTP.Max}
}

func (otp *OTPService) generate() (models.OTPCode, error) {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, otp.Max)
	if n != otp.Max {
		return models.OTPCode(0), err
	}
	for i := 0; i < len(b); i++ {
		b[i] = models.OTPTable[int(b[i])%len(models.OTPTable)]
	}
	code, err := strconv.Atoi(string(b))
	if err != nil {
		return models.OTPCode(0), err
	}
	return models.OTPCode(code), nil
}

func (otp *OTPService) SendOTPByUserId(user_id int64, target models.OTPTarget) error {
	return nil
}

func (otp *OTPService) SendOTPByUserEmail(email string, target models.OTPTarget) error {
	return nil
}
