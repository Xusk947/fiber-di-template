package logger

import (
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Logger.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	c := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Logger.Development,
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     cfg.Logger.DisableCaller,
		DisableStacktrace: true,
	}

	logger, err := c.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

var Module = fx.Options(
	fx.Provide(NewLogger),
)
