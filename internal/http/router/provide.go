package router

import "go.uber.org/fx"

type RouteConstructor[P any, R any] func(P) R

type RouteOptions[P HasRouteParamsInjected, R any] struct {
	Constructor RouteConstructor[P, R]
	Acc         AccessType
}

func ProvideRoute[P HasRouteParamsInjected, R any](opts RouteOptions[P, R]) any {
	acc := opts.Acc
	if acc == "" {
		acc = RoutePublic
	}
	return fx.Annotate(
		opts.Constructor,
		fx.As(new(R)),
		fx.ResultTags(string(acc)),
	)
}
