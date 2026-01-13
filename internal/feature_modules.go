package internal

import (
	"rest-fiber/internal/auth"
	"rest-fiber/internal/user"

	"go.uber.org/fx"
)

var FeatureModules = fx.Options(
	user.Module,
	auth.Module,
	// add more modules in here
)
