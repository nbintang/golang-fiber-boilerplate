package auth

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/http/router" 
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		NewAuthService,
		NewAuthHandler,
		router.ProvideRoute[AuthRouteParams, router.Route](
			NewAuthRoute,
			enums.RoutePublic,
		),
	),
)
