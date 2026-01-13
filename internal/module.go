package internal

import "go.uber.org/fx"

var Module = fx.Module(
	"app",
	fx.Provide(NewBootstrap),
	fx.Invoke(
		RegisterRoutes,
		RegisterLifecycle,
	),
	FeatureModules,
)
