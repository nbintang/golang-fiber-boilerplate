package middleware

import (
	"rest-fiber/internal/identity"
	"rest-fiber/internal/infra/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AccessNotBlacklisted(cacheService cache.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := identity.CurrentUser(c)
		if err != nil {
			return err
		}
		key := "blacklist_access:" + user.JTI
		ctx := c.UserContext()
		_, err = cacheService.Get(ctx, key)
		if err == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		if err != redis.Nil {
			return fiber.NewError(fiber.StatusInternalServerError, "redis error")
		}
		return c.Next()
	}
}
