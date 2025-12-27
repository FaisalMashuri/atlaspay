package repository

import (
	"auth_service/infrastructure/database"
	"auth_service/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) IRepository {
	return &repository{
		db: db,
	}
}

func (repo *repository) CreateUser(ctx context.Context, exec database.Executor, user *model.User) error {
	query := `INSERT INTO USERS (created_at, updated_at, user_ref, email, password_hash, status) 
		VALUES 
		(:created_at, :updated_at, :user_ref, :email, :password_hash, :status);`

	_, err := exec.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}
