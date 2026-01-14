package auth

import (
	"go.uber.org/fx"
	"rest-fiber/internal/http/router"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		router.ProvideRoute[AuthRouteParams, router.Route](
			router.RouteOptions[AuthRouteParams, router.Route]{
				Constructor: NewAuthRoute,
				Acc:         router.RouteProtected,
			},
		),
	),
)
