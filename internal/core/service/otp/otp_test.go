package otp

import (
	"dashboard/internal/common/service/config"
	"testing"
)

func TestGenerate(t *testing.T) {
	cfg := config.AppConfig{OTP: config.OTP{Max: 6}}
	otp := NewOTPService(&cfg)
	code, err := otp.generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !code.IsValid() {
		t.Fatalf("code is invalid: %d", code)
	}
}
