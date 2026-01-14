package internal

import (
	"rest-fiber/config"
	"rest-fiber/internal/http"
	"rest-fiber/internal/http/middleware"
	"rest-fiber/internal/infra/cache"
	"rest-fiber/internal/infra/infraapp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Bootstrap struct {
	*fiber.App
	PublicRoute    fiber.Router
	ProtectedRoute fiber.Router
	Env            config.Env
	Logger         *infraapp.AppLogger
}

func NewBootstrap(env config.Env, logger *infraapp.AppLogger, cacheService cache.Service) *Bootstrap {
	app := fiber.New(fiber.Config{
		ErrorHandler: http.DefaultErrorHandler,
		AppName:      "Fiber Rest API",
	})
	app.Use(middleware.LoggerRequest(logger))
	app.Use(middleware.RequestMeta())
	app.Use(cors.New(cors.ConfigDefault))
	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Welcome to API v1")
	})

	protected := v1.Group("/protected")
	protected.Use(
		middleware.AccessToken(env),
		middleware.AccessCurrentUser(),
		middleware.AccessNotBlacklisted(cacheService),
	)

	return &Bootstrap{
		App:            app,
		PublicRoute:    v1,
		ProtectedRoute: protected,
		Env:            env,
		Logger:         logger,
	}
}
