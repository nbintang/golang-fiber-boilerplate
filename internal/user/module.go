package user

import (
	"rest-fiber/internal/http/router"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewUserRepository,
		NewUserService,
		NewUserHandler,
		router.ProvideRoute[UserRouteParams, router.ProtectedRoute](
			router.ProvideRouteOptions[UserRouteParams, router.ProtectedRoute]{
				Constructor: NewUserRoute,
				Acc:         router.RouteProtected,
			},
		),
	),
)
