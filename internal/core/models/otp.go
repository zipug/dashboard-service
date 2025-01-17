package models

type OTPCode int

type OTPTarget string

const (
	VERIFICATION   OTPTarget = "verify"
	AUTHENTICATION OTPTarget = "auth"
)

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
