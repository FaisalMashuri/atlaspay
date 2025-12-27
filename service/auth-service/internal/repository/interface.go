package repository

import (
	"auth_service/infrastructure/database"
	"auth_service/internal/model"
	"context"
)

type IRepository interface {
	CreateUser(ctx context.Context, exec database.Executor, user *model.User) error
}
