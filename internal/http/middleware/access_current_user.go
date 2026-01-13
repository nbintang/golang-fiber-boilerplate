package middleware

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/identity"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AccessCurrentUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals(enums.AccessAuthKey).(*jwt.Token)
		if !ok || token == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		id, _ := claims["id"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)
		jti, _ := claims["jti"].(string)
		if id == "" ||
			email == "" ||
			role == "" ||
			jti == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		c.Locals(enums.CurrentUserKey, &identity.AuthClaims{
			ID:    id,
			Email: email,
			Role:  enums.EUserRoleType(role),
			JTI:   jti,
		})
		return c.Next()
	}
}
