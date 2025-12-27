package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Application AppConfig      `mapstructure:",squash"`
	Database    DatabaseConfig `mapstructure:",squash"`
}

func LoadConfig() (Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("failed load .env:", err)
		}
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
	v.SetDefault("database_driver", "postgres")
	v.SetDefault("database_host", "")
	v.SetDefault("database_port", 0)
	v.SetDefault("database_name", "")
	v.SetDefault("database_user", "")
	v.SetDefault("database_password", "")
	v.SetDefault("database_ssl_mode", "disable")
	v.SetDefault("database_timezone", "UTC")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	
	validate(&cfg)
	return cfg, nil
}

func validate(cfg *Config) {
	if cfg.Database.Host == "" || cfg.Database.Port == 0 || cfg.Database.User == "" || cfg.Database.Password == "" {
		panic("Config database is required")
	}
}

func (c Config) Print() {
	fmt.Println("CONFIG LOADED")
	fmt.Println("=============")

	fmt.Println("[APPLICATION]")
	fmt.Printf("ENV              : %s\n", c.Application.Environment)
	fmt.Printf("PORT             : %d\n", c.Application.Port)
	fmt.Printf("LOG LEVEL        : %s\n", c.Application.LogLevel)
	fmt.Printf("DEBUG            : %v\n", c.Application.EnableDebug)
	fmt.Printf("SHUTDOWN TIMEOUT : %s\n", c.Application.ShutdownTimeout)
	fmt.Printf("READ TIMEOUT     : %s\n", c.Application.ReadTimeout)
	fmt.Printf("WRITE TIMEOUT    : %s\n", c.Application.WriteTimeout)
	fmt.Printf("IDLE TIMEOUT     : %s\n", c.Application.IdleTimeout)

	fmt.Println()
	fmt.Println("[DATABASE]")
	fmt.Printf("DRIVER           : %s\n", c.Database.Driver)
	fmt.Printf("HOST             : %s\n", c.Database.Host)
	fmt.Printf("PORT             : %d\n", c.Database.Port)
	fmt.Printf("NAME             : %s\n", c.Database.Name)
	fmt.Printf("USER             : %s\n", c.Database.User)
	fmt.Printf("PASSWORD         : %s\n", mask(c.Database.Password))
	fmt.Printf("SSL MODE         : %s\n", c.Database.SSLMode)
	fmt.Printf("TIMEZONE         : %s\n", c.Database.TimeZone)
	fmt.Printf("MAX OPEN CONNS   : %d\n", c.Database.MaxOpenConns)
	fmt.Printf("MAX IDLE CONNS   : %d\n", c.Database.MaxIdleConns)
	fmt.Printf("CONN MAX LIFE    : %s\n", c.Database.ConnMaxLifetime)
	fmt.Printf("CONN IDLE TIME   : %s\n", c.Database.ConnMaxIdleTime)
}
func mask(s string) string {
	if s == "" {
		return "(empty)"
	}
	return "******"
}
