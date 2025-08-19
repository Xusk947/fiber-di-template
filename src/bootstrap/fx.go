package bootstrap

import (
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/logger"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/probe"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/server"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log.Named("fx").WithOptions(zap.IncreaseLevel(zap.ErrorLevel))}
	}),
	//postgres.Module,
	//clickhouse.Module,
	//redis.Module,
	//kafka.Module,
	probe.Module,
	server.Module,
	internal.Module,
)
