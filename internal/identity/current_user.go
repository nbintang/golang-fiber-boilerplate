package identity

import (
	"rest-fiber/internal/apperr"
	"rest-fiber/internal/enums"

	"github.com/gofiber/fiber/v2"
)

func CurrentUser[T AuthClaims](c *fiber.Ctx) (*T, error) {
	v := c.Locals(enums.CurrentUserKey)
	user, ok := v.(*T)
	if !ok || user == nil {
		return nil, apperr.Unauthorized(apperr.CodeUnauthorized, "Unauthorized", nil)
	}
	return user, nil
}
