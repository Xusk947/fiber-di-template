package internal

import (
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal/routes"
	"go.uber.org/fx"
)

var Module = fx.Module("internal", fx.Options(
	/* Add new modules here */

	// Do not remove this line, always leave it at the end
	routes.Module,
))
