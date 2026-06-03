package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(debug string, env string) (*zap.Logger, error) {
	var cfg zap.Config
	var err error

	if env == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.Level, err = zap.ParseAtomicLevel(debug)
	if err != nil {
		return nil, err
	}

	return cfg.Build()
}
