package service

import (
	"auth_service/infrastructure/database"
	dbmocks "auth_service/infrastructure/database/mocks"
	"auth_service/internal/model"
	"auth_service/internal/model/dto/requests"
	repomocks "auth_service/internal/repository/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockTxSuccess(tx *dbmocks.MockITxManager) {
	tx.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			ctx context.Context,
			fn func(database.Executor) error,
		) error {
			return fn(nil)
		})
}

func mockTxError(tx *dbmocks.MockITxManager, err error) {
	tx.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		Return(err)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		req         requests.RegisterRequest
		mockSetup   func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository)
		expectError bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			req: requests.RegisterRequest{
				Email:    "test@mail.com",
				Password: "secret",
			},
			mockSetup: func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository) {
				mockTxSuccess(tx)

				repo.EXPECT().
					CreateUser(
						gomock.Any(),
						gomock.Any(),
						gomock.AssignableToTypeOf(&model.User{}),
					).
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "mapping error",
			ctx:  context.Background(),
			req:  requests.RegisterRequest{}, // invalid
			mockSetup: func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository) {
				// nothing should be called
			},
			expectError: true,
		},
		{
			name: "repo error",
			ctx:  context.Background(),
			req: requests.RegisterRequest{
				Email:    "fail@mail.com",
				Password: "secret",
			},
			mockSetup: func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository) {
				mockTxSuccess(tx)

				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("repo error"))
			},
			expectError: true,
		},
		{
			name: "tx error",
			ctx:  context.Background(),
			req: requests.RegisterRequest{
				Email:    "tx@mail.com",
				Password: "secret",
			},
			mockSetup: func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository) {
				mockTxError(tx, errors.New("tx failed"))
			},
			expectError: true,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			req: requests.RegisterRequest{
				Email:    "ctx@mail.com",
				Password: "secret",
			},
			mockSetup: func(tx *dbmocks.MockITxManager, repo *repomocks.MockIRepository) {
				mockTxError(tx, context.Canceled)
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tx := dbmocks.NewMockITxManager(ctrl)
			repo := repomocks.NewMockIRepository(ctrl)

			if tt.mockSetup != nil {
				tt.mockSetup(tx, repo)
			}

			svc := NewService(repo, tx)

			err := svc.Register(tt.ctx, tt.req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
