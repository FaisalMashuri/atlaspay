package repository

import (
	"auth_service/infrastructure/database"
	"context"
	"github.com/jmoiron/sqlx"
	"os/user"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) IRepository {
	return &repository{
		db: db,
	}
}

func (repo *repository) CreateUser(ctx context.Context, exec database.Executor, user *user.User) error {
	query := `INSERT INTO USERS (created_at, updated_at, user_ref, email, password_hash, status) 
		VALUES 
		(:created_at, :updated_at, :user_ref, :email, :password, :status);`

	_, err := exec.ExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}
