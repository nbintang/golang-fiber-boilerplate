package identity

import (
	"rest-fiber/internal/enums"

	"github.com/gofiber/fiber/v2"
)

func CurrentUser[T AuthClaims](c *fiber.Ctx) (*T, error) {
	v := c.Locals(enums.CurrentUserKey)
	user, ok := v.(*T)
	if !ok || user == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return user, nil
}
