package internal

import (
	"rest-fiber/config"
	"rest-fiber/internal/http"
	"rest-fiber/internal/http/middleware"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/rediscache"

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

func NewBootstrap(env config.Env, logger *infraapp.AppLogger, redisService rediscache.Service) *Bootstrap {
	app := fiber.New(fiber.Config{
		ErrorHandler: http.DefaultErrorHandler,
		AppName:      "Fiber Rest API",
	})
	app.Use(middleware.LoggerRequest(logger))
	app.Use(middleware.RequestMeta())
	app.Use(cors.New(cors.ConfigDefault))
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Wellcome to API")
	})

	protected := api.Group("/protected")
	// Protected Routes Provider
	protected.Use(
		middleware.AccessToken(env),
		middleware.AccessCurrentUser(),
		middleware.AccessNotBlacklisted(redisService),
	)

	return &Bootstrap{
		App:            app,
		PublicRoute:    api,
		ProtectedRoute: protected,
		Env:            env,
		Logger:         logger,
	}
}
