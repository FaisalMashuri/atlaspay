package requests

import (
	"auth_service/internal/model"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
}

func (req *RegisterRequest) ToModelNewUser() (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	now := time.Now()

	return model.User{
		Email:        req.Email,
		CreatedAt:    now,
		UpdatedAt:    now,
		Status:       "ACTIVE",
		PasswordHash: string(passwordHash),
		UserRef:      ulid.Make().String(),
	}, nil
}
