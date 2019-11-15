package otp

import (
	"github.com/pquerna/otp/totp"
	"time"
)

func GetTOTPPassCode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
