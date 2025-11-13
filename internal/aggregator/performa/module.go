package performa

import "go.uber.org/fx"

var Module = fx.Module(
	"performaAgg",
	fx.Provide(NewPerformaAggregator),
)
