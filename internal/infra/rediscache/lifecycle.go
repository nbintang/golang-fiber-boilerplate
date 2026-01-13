package rediscache

import (
	"context"

	"go.uber.org/fx"
)

func RegisterLifecycle(
	lc fx.Lifecycle,
	redisSvc Service,
) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return redisSvc.Close()
		},
	})
}
