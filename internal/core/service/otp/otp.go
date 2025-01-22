package otp

import (
	"context"
	"crypto/rand"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"io"
	"strconv"
)

type OTPService struct {
	Max  int
	repo ports.OTPRepository
}

func NewOTPService(cfg *config.AppConfig, repo ports.OTPRepository) *OTPService {
	return &OTPService{Max: cfg.OTP.Max, repo: repo}
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

func (otp *OTPService) SendOTP(ctx context.Context, user_id int64, email string, username string) error {
	code, err := otp.generate()
	if err != nil {
		return err
	}
	return otp.repo.SaveUserOTP(ctx, user_id, email, username, code)
}

func (otp *OTPService) GetOTP(ctx context.Context, user_id int64) (models.OTPCode, error) {
	return otp.repo.GetUserOTP(ctx, user_id)
}
