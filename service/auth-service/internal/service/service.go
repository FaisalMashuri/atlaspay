package service

import (
	"auth_service/infrastructure/database"
	"auth_service/internal/model/dto/requests"
	"auth_service/internal/repository"
	"context"
)

type service struct {
	tx   database.TxManager
	repo repository.IRepository
}

func NewService(repo repository.IRepository, tx database.TxManager) IService {
	return &service{
		repo: repo,
		tx:   tx,
	}
}

func (s *service) Register(ctx context.Context, reqData requests.RegisterRequest) error {
	user, errMapping := reqData.ToModelNewUser()
	if errMapping != nil {
		return errMapping
	}
	err := s.tx.WithTransaction(ctx, func(exec database.Executor) error {
		err := s.repo.CreateUser(ctx, exec, &user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
