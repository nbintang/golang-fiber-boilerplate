package http

import (
	"errors"

	"rest-fiber/internal/apperr"
	"rest-fiber/pkg/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	var detail any = nil
	var code any = nil

	switch {
	case errors.As(err, new(*apperr.AppError)):
		var ae *apperr.AppError
		_ = errors.As(err, &ae)

		statusCode = ae.Status
		msg = ae.Message
		code = string(ae.Code)

	case errors.As(err, new(validator.ValidationErrors)):
		var ve validator.ValidationErrors
		_ = errors.As(err, &ve)

		out := make([]fiber.Map, 0, len(ve))
		for _, fe := range ve {
			out = append(out, fiber.Map{
				"field": fe.Field(),
				"tag":   fe.Tag(),
			})
		}

		statusCode = fiber.StatusUnprocessableEntity
		msg = "Validation Error"
		code = "VALIDATION_ERROR"
		detail = out

	case errors.As(err, new(*fiber.Error)):
		var fe *fiber.Error
		_ = errors.As(err, &fe)

		statusCode = fe.Code
		msg = fe.Message
		code = "FIBER_ERROR"

	default:
		code = "INTERNAL"
	}

	return c.Status(statusCode).JSON(httpx.NewHttpResponse(
		statusCode,
		msg,
		fiber.Map{
			"code":  code,
			"error": detail,
			"meta": fiber.Map{
				"method":   c.Locals("method"),
				"path":     c.Locals("path"),
				"endpoint": c.Locals("endpoint"),
				"status":   statusCode,
				"latency":  c.Locals("latency"),
				"ip":       c.Locals("ip"),
			},
		},
	))
}
