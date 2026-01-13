package main

import (
	"rest-fiber/config"
	app "rest-fiber/internal"
	"rest-fiber/internal/infra"
	"rest-fiber/pkg/env"

	"go.uber.org/fx"
)

func main() {
	env.Load()
	fx.New(
		config.Module,
		infra.Module,
		app.Module,
	).Run()
}
 