package postgres

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

var Module = fx.Module("postgres",
	fx.Provide(
		// Provide the main Postgres struct
		NewPostgres,
		// Also provide the underlying *sqlx.DB client for repositories
		func(p *Postgres) *sqlx.DB {
			return p.Client
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, db *Postgres) {
		lc.Append(fx.Hook{
			OnStart: db.onStart,
			OnStop:  db.onStop,
		})
	}),
)
