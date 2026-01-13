package http

import (
	"rest-fiber/pkg/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	msg := "Internal Server Error"
	var data any = nil

	if ve, ok := err.(validator.ValidationErrors); ok {
		out := make([]fiber.Map, 0, len(ve))
		for _, fe := range ve {
			out = append(out, fiber.Map{
				"field": fe.Field(),
				"tag":   fe.Tag(),
			})
		}
		statusCode = fiber.StatusBadRequest
		data = out
	}
 
	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
		msg = e.Message
	}

	meta := fiber.Map{
		"method":   c.Locals("method"),
		"path":     c.Locals("path"),
		"endpoint": c.Locals("endpoint"),
		"status":   statusCode,
		"latency":  c.Locals("latency"),
		"ip":       c.Locals("ip"),
	}

	return c.Status(statusCode).JSON(httpx.NewHttpResponse(
		statusCode,
		msg,
		fiber.Map{
			"error": data,
			"meta":  meta,
		},
	))
}
