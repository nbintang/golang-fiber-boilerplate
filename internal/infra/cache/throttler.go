package cache

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Throttle(params ThrottleParams) fiber.Handler {
	exp := params.Expiration;
	if exp <= 0 {
		exp = 1 * time.Minute
	}

	prefix := params.Prefix;
	if prefix == "" {
		prefix = "route"
	}

	config := limiter.Config{
		Max:        params.MaxLimit,
		Expiration: exp,
		KeyGenerator: func(c *fiber.Ctx) string {
			return prefix + ":" + c.Path() + ":" + c.IP()
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too many requests")
		},
	}
	if params.Storage != nil {
		config.Storage = params.Storage
	}
	return limiter.New(config)
}
