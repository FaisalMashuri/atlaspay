package repository

import (
	"auth_service/infrastructure/database"
	"context"
	"os/user"
)

type IRepository interface {
	CreateUser(ctx context.Context, exec database.Executor, user *user.User) error
}
