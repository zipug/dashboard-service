package models

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidOTPCode = errors.New("invalid otp code")
	ErrCodeNotEqualTo = errors.New("otp code not equal to")
)

type OTPCode int

var OTPTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func (code OTPCode) len() int {
	if code == 0 {
		return 1
	}
	count := 0
	for code != 0 {
		code /= 10
		count++
	}
	return count
}

func (code OTPCode) IsValid() bool {
	if code < 0 {
		return false
	}
	return code.len() == 6
}

func (code OTPCode) ToStr() (string, error) {
	otp := strconv.Itoa(int(code))
	if otp == "" {
		return otp, ErrInvalidOTPCode
	}
	return otp, nil
}

func (code OTPCode) Equals(new_code OTPCode) bool {
	return code == new_code
}

func StrToOTPCode(otp string) (OTPCode, error) {
	code, err := strconv.Atoi(otp)
	if err != nil {
		return OTPCode(0), err
	}
	return OTPCode(code), nil
}
