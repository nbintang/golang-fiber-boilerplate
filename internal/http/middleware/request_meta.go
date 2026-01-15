package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestMeta() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		
		c.Locals("method", c.Method())
		c.Locals("path", c.OriginalURL())
		if c.Route() != nil {
			c.Locals("endpoint", c.Route().Path)
		}
		c.Locals("status", c.Response().StatusCode())
		c.Locals("latency", time.Since(start).String())
		c.Locals("ip", c.IP())

		return err
	}
}