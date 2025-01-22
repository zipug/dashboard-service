package otp

import (
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/ports"
	"testing"
)

func TestGenerate(t *testing.T) {
	cfg := config.AppConfig{OTP: config.OTP{Max: 6}}
	var mockRepo ports.OTPRepository
	otp := NewOTPService(&cfg, mockRepo)
	code, err := otp.generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !code.IsValid() {
		t.Fatalf("code is invalid: %d", code)
	}
}
