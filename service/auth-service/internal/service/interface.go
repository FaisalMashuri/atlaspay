package service

import (
	"auth_service/internal/model/dto/requests"
	"context"
)

type IService interface {
	Register(ctx context.Context, reqData requests.RegisterRequest) error
}
