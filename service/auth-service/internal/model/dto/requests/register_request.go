package requests

import (
	"auth_service/internal/model"
	"errors"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=8,max=20"`
}

func (req *RegisterRequest) ToModelNewUser() (model.User, error) {
	if req.Email == "" || req.Password == "" {
		return model.User{}, errors.New("email or password is empty")
	}
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
