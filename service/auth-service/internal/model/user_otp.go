package model

import (
	"auth_service/utils/constant"
	"time"
)

type UserOTP struct {
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`

	OTPCodeHash string              `db:"otp_code_hash"`
	Purpose     constant.OTPPurpose `db:"purpose"`

	UsedAt *time.Time `db:"used_at"`

	UserID int64 `db:"user_id"`
	ID     int64 `db:"id"`
}
