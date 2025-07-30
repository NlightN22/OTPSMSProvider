package service

import "time"

// OTPService defines business logic for TOTP.
type OTPService interface {
	GenerateCode(key string) (code string, err error)
	ValidateCode(key, code string) bool
	CanSend(key string) (bool, time.Duration)
}
