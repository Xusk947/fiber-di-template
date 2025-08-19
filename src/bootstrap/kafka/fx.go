package kafka

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
)

var Module = fx.Module("kafka",
	fx.Provide(
		// Provide the main Kafka client
		NewKafka,
	),
	fx.Invoke(func(lc fx.Lifecycle, kafka *Kafka) {
		lc.Append(fx.Hook{
			OnStart: kafka.onStart,
			OnStop:  kafka.onStop,
		})
	}),
)
