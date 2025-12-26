package model

import (
	"auth_service/utils/constant"
	"time"
)

type User struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Email        string              `db:"email"`
	PasswordHash string              `db:"password_hash"`
	UserRef      string              `db:"user_ref"`
	Status       constant.UserStatus `db:"status"`

	ID int64 `db:"id"`
}
