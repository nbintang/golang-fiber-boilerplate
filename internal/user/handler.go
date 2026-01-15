package user

import (
	"rest-fiber/internal/apperr"
	"rest-fiber/internal/identity"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/validator"
	"rest-fiber/pkg/httpx"
	"rest-fiber/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService UserService
	logger      *infraapp.AppLogger
	validator   validator.Service
}

func NewUserHandler(userService UserService, logger *infraapp.AppLogger, validator validator.Service) UserHandler {
	return &userHandlerImpl{userService, logger, validator}
}

func (h *userHandlerImpl) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var query pagination.Query
	if err := c.QueryParser(&query); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Invalid Request", err)
	}

	query = query.Normalize(10, 100)

	data, total, err := h.userService.FindAllUsers(ctx, query.Page, query.Limit, query.Offset())
	if err != nil {
		return err
	}

	meta := pagination.NewMeta(query.Page, query.Limit, total)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpPaginationResponse[[]UserResponseDTO](
		fiber.StatusOK,
		"Success",
		data,
		meta,
	))
}

func (h *userHandlerImpl) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Invalid user ID", err)
	}
	ctx := c.UserContext()
	data, err := h.userService.FindUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", data))
}

func (h *userHandlerImpl) GetCurrentUserProfile(c *fiber.Ctx) error {
	currentUser, err := identity.CurrentUser(c)
	if err != nil {
		return err
	}
	ctx := c.UserContext()
	h.logger.Infof("user Id :%s", currentUser.ID)
	data, err := h.userService.FindUserByID(ctx, currentUser.ID)
	if err != nil {
		return  err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(fiber.StatusOK, "Success", data))
}

func (h *userHandlerImpl) UpdateCurrentUser(c *fiber.Ctx) error {
	currentUser, err := identity.CurrentUser(c)
	if err != nil {
		return err
	}
	var body UserUpdateDTO
	if err := c.BodyParser(&body); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Invalid Request", err)
	}

	if err := h.validator.Struct(body); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Validation Error", err)
	}

	if err := h.userService.UpdateProfile(c.UserContext(), currentUser.ID, body); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(
		httpx.NewHttpResponse[any](fiber.StatusOK, "User Updated Successfully", nil),
	)
}
