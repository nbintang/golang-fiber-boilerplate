package internal

import (
	"context"

	"go.uber.org/fx"
)

func RegisterLifecycle(lc fx.Lifecycle, a *Bootstrap) {
	addr := a.Env.AppAddr
	if addr == "" {
		addr = ":8080"
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				a.Logger.Printf("Fiber listening on %s", addr)
				if err := a.Listen(addr); err != nil {
					a.Logger.Printf("Fiber stopped: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return a.Shutdown()
		},
	})
}
