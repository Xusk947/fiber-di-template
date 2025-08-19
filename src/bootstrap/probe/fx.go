package probe

import (
	"go.uber.org/fx"
)

var Module = fx.Module("probe",
	fx.Provide(
		NewProbeServer,
	),
	fx.Invoke(func(lc fx.Lifecycle, probe *ProbeServer) {
		lc.Append(fx.Hook{
			OnStart: probe.Start,
			OnStop:  probe.Stop,
		})
	}),
)
