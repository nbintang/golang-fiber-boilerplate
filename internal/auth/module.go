package auth

import (
	"rest-fiber/internal/http/router"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		router.ProvideRoute[AuthRouteParams, router.Route](
			router.ProvideRouteOptions[AuthRouteParams, router.Route]{
				Constructor: NewAuthRoute,
				Acc:         router.RouteProtected,
			},
		),
	),
)
