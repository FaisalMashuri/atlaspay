package database

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ITxManager interface {
	WithTransaction(
		ctx context.Context,
		fn func(exec Executor) error,
	) error
}

type TxManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) ITxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTransaction(ctx context.Context, fn func(exec Executor) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
