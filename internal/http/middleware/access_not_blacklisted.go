package middleware

import (
	"rest-fiber/internal/apperr"
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
		if _, err = cacheService.Get(ctx, key); err == nil {
			return apperr.Forbidden(apperr.CodeForbidden, "Access Denied", err)
		}
		if err != redis.Nil {
			return apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
		}
		return c.Next()
	}
}
