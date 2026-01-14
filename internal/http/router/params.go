package router

import "go.uber.org/fx"

type InjectRouteParams struct {
	fx.In
}

func (InjectRouteParams) hasRouteParams() {}

type HasRouteParamsInjected interface {
	hasRouteParams()
}
