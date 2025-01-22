package models

type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

type OTPMessageType string

const (
	Verify OTPMessageType = "verify"
	Login  OTPMessageType = "login"
)

type OTPPayload struct {
	Type     OTPMessageType `json:"type"`
	UserID   int64          `json:"user_id"`
	UserName string         `json:"username"`
	Email    string         `json:"email"`
	Code     string         `json:"code"`
}
