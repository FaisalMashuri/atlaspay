package repository

import (
	"auth_service/infrastructure/database/mocks"
	"auth_service/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	now := time.Now()
	testTable := []struct {
		name        string
		ctx         context.Context
		user        *model.User
		mockSetup   func(exec *mocks.MockExecutor, ctx context.Context, user *model.User)
		expectError bool
	}{
		{
			name: "success insert user",
			ctx:  context.Background(),
			user: &model.User{
				UserRef:      "USR-001",
				Email:        "test@mail.com",
				PasswordHash: "hashed",
				Status:       "ACTIVE",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			mockSetup: func(exec *mocks.MockExecutor, ctx context.Context, user *model.User) {
				exec.EXPECT().
					NamedExecContext(ctx, gomock.Any(), user).
					Return(sql.Result(nil), nil)
			},
			expectError: false,
		},
		{
			name: "database error",
			ctx:  context.Background(),
			user: &model.User{
				UserRef: "USR-002",
				Email:   "fail@mail.com",
			},
			mockSetup: func(exec *mocks.MockExecutor, ctx context.Context, user *model.User) {
				exec.EXPECT().
					NamedExecContext(ctx, gomock.Any(), user).
					Return(nil, errors.New("db error"))
			},
			expectError: true,
		},
		{
			name: "empty email",
			ctx:  context.Background(),
			user: &model.User{
				UserRef: "USR-003",
				Email:   "",
			},
			mockSetup: func(exec *mocks.MockExecutor, ctx context.Context, user *model.User) {
				exec.EXPECT().
					NamedExecContext(ctx, gomock.Any(), user).
					Return(sql.Result(nil), nil)
			},
			expectError: false, // repository tidak validasi field
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			user: &model.User{
				UserRef: "USR-004",
				Email:   "cancel@mail.com",
			},
			mockSetup: func(exec *mocks.MockExecutor, ctx context.Context, user *model.User) {
				exec.EXPECT().
					NamedExecContext(ctx, gomock.Any(), user).
					Return(nil, context.Canceled)
			},
			expectError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			exec := mocks.NewMockExecutor(ctrl)
			repo := NewRepository(nil)
			if tt.mockSetup != nil {
				tt.mockSetup(exec, tt.ctx, tt.user)
			}
			var err error
			assert.NotPanics(t, func() {
				err = repo.CreateUser(tt.ctx, exec, tt.user)
			})
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
