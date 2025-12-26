package model

import "time"

type RefreshToken struct {
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`

	TokenHash string `db:"token_hash"`

	RevokedAt *time.Time `db:"revoked_at"`

	UserID int64 `db:"user_id"`
	ID     int64 `db:"id"`
}
