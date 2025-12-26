package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	v := viper.New()

	// ENV only
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ---- SAFE DEFAULTS ----
	v.SetDefault("server_port", 8080)
	v.SetDefault("server_read_timeout", "5s")
	v.SetDefault("server_write_timeout", "10s")
	v.SetDefault("server_idle_timeout", "60s")

	v.SetDefault("database_max_open_conns", 25)
	v.SetDefault("database_max_idle_conns", 10)
	v.SetDefault("database_conn_max_lifetime", "30m")
	v.SetDefault("database_conn_max_idle_time", "5m")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate(&cfg)
	return &cfg, nil
}

func validate(cfg *Config) {
	if cfg.Database.DSN == "" {
		panic("DATABASE_DSN is required")
	}
}
