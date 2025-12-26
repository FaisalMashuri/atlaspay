package database

import (
	"auth_service/config"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

var (
	instance *sqlx.DB
	once     sync.Once
	initErr  error
)

func ConnectDatabase(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	once.Do(func() {
		dsn := buildPostgresDSN(cfg)

		instance, initErr = sqlx.Connect("postgres", dsn)
		if initErr != nil {
			return
		}

		instance.SetMaxOpenConns(cfg.MaxOpenConns)
		instance.SetMaxIdleConns(cfg.MaxIdleConns)
		instance.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
		instance.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := instance.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping database due to %w", err)
			_ = instance.Close()
			instance = nil
			return
		}
	})
	return instance, initErr
}

func buildPostgresDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s connect_timeout=5",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}

func Ping(ctx context.Context) error {
	if instance == nil {
		return nil
	}
	return instance.PingContext(ctx)
}
