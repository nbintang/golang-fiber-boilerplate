package router

import (
	"rest-fiber/internal/enums"

	"go.uber.org/fx"
)

type RouteConstructor[P any, R any] func(P) R

func ProvideRoute[P any, R any](
	routeConstructor RouteConstructor[P, R],
	acc enums.EAccessType,
) any {
	return fx.Annotate(
		routeConstructor,
		fx.As(new(R)),
		fx.ResultTags(string(acc)),
	)
}
