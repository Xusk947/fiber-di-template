package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/fx"
)

var Module = fx.Module("clickhouse",
	fx.Provide(
		// Provide the main ClickHouse struct
		NewClickHouse,
		// Also provide the underlying driver.Conn client for repositories
		func(c *ClickHouse) driver.Conn {
			return c.Client
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, db *ClickHouse) {
		lc.Append(fx.Hook{
			OnStart: db.onStart,
			OnStop:  db.onStop,
		})
	}),
)
