package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	NATS     NATSConfig
	Log      LogConfig
}

type AppConfig struct {
	Env  string `env:"APP_ENV" env-default:"development"`
	Name string `env:"APP_NAME" env-default:"analytics-service"`
}

type PostgresConfig struct {
	DSN      string `env:"POSTGRES_DSN" env-required:"true"`
	MaxConns int    `env:"POSTGRES_MAX_CONNS" env-default:"20"`
	MinConns int    `env:"POSTGRES_MIN_CONNS" env-default:"2"`
}

type NATSConfig struct {
	URL            string `env:"NATS_URL" env-default:"nats://localhost:4222"`
	MaxReconnects  int    `env:"NATS_MAX_RECONNECTS" env-default:"10"`
	ReconnectWaitS int    `env:"NATS_RECONNECT_WAIT_S" env-default:"2"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
