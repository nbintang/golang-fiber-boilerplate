package infra

import (
	"rest-fiber/internal/infra/cache"
	"rest-fiber/internal/infra/database"
	"rest-fiber/internal/infra/email"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/token"
	"rest-fiber/internal/infra/validator"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infra",
	fx.Provide(
		database.NewService,
		cache.NewService,
	),
	fx.Provide(
		validator.NewService,
		token.NewService,
		email.NewService,
	),
	fx.Provide(
		infraapp.NewLogger,
		database.NewLogger,
	),
	fx.Invoke(
		cache.RegisterLifecycle,
		database.RegisterLifecycle,
	),
)
