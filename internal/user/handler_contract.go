package user

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetCurrentUserProfile(c *fiber.Ctx) error
	UpdateCurrentUser(c *fiber.Ctx) error
}
