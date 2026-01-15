package middleware

import (
	"rest-fiber/internal/apperr"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/identity"

	"github.com/gofiber/fiber/v2"
)

func AllowRoleAccess(roles ...enums.EUserRoleType) fiber.Handler {
	roleSet := make(map[enums.EUserRoleType]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		user, err := identity.CurrentUser(c)
		if err != nil {
			return err
		}
		if _, ok := roleSet[user.Role]; !ok {
			return apperr.Forbidden(apperr.CodeForbidden, "Forbidden", nil)
		}
		return c.Next()
	}
}
