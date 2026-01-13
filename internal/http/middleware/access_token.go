package middleware

import (
	"rest-fiber/config"
	"rest-fiber/internal/enums"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func AccessToken(env config.Env) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(env.JWTAccessSecret)},
		ContextKey: enums.AccessAuthKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message":     "unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"error":       err.Error(),
			})
		},
	})
}
