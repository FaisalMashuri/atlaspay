package config

import "time"

type AppConfig struct {
	Port            int           `mapstructure:"server_port"`
	Environment     string        `mapstructure:"app_env"` // dev | staging | prod
	LogLevel        string        `mapstructure:"app_log_level"`
	ShutdownTimeout time.Duration `mapstructure:"app_shutdown_timeout"`
	EnableDebug     bool          `mapstructure:"app_enable_debug"`
	ReadTimeout     time.Duration `mapstructure:"server_read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"server_write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"server_idle_timeout"`
}
