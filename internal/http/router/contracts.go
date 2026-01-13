package router

import "github.com/gofiber/fiber/v2"

type Route interface {
	RegisterRoute(route fiber.Router)
}

type ProtectedRoute interface {
	RegisterProtectedRoute(route fiber.Router)
}