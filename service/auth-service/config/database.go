package config

import "time"

type DatabaseConfig struct {
	Driver string `mapstructure:"database_driver"`

	Host string `mapstructure:"database_host"`
	Port int    `mapstructure:"database_port"`

	Name     string `mapstructure:"database_name"`
	User     string `mapstructure:"database_user"`
	Password string `mapstructure:"database_password"`

	SSLMode  string `mapstructure:"database_ssl_mode"` // disable | require | verify-full
	TimeZone string `mapstructure:"database_timezone"`

	MaxOpenConns    int           `mapstructure:"database_max_open_conns"`
	MaxIdleConns    int           `mapstructure:"database_max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"database_conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"database_conn_max_idle_time"`
}
