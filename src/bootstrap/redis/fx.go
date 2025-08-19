package redis

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var Module = fx.Module("redis",
	fx.Provide(
		// Provide the main Redis struct
		NewRedis,
		// Also provide the underlying *redis.Client for services
		func(r *Redis) *redis.Client {
			return r.Client
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, redis *Redis) {
		lc.Append(fx.Hook{
			OnStart: redis.onStart,
			OnStop:  redis.onStop,
		})
	}),
)
