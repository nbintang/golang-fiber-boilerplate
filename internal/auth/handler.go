package auth

import (
	"errors"
	"rest-fiber/config"
	"rest-fiber/internal/apperr"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/validator"
	"rest-fiber/pkg/httpx"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authHandlerImpl struct {
	authService AuthService
	validate    validator.Service
	env         config.Env
	logger      *infraapp.AppLogger
}

func NewAuthHandler(authService AuthService, validate validator.Service, env config.Env, logger *infraapp.AppLogger) AuthHandler {
	return &authHandlerImpl{authService, validate, env, logger}
}

func (h *authHandlerImpl) Register(c *fiber.Ctx) error {
	var body RegisterRequestDTO
	ctx := c.UserContext()

	if err := c.BodyParser(&body); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Invalid Request", err)
	}
	if err := h.validate.Struct(body); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Validation Error", err)
	}
	if err := h.authService.Register(ctx, &body); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(httpx.NewHttpResponse[any](
		fiber.StatusCreated,
		"Success! please check your email",
		nil,
	))
}

func (h *authHandlerImpl) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	ctx := c.UserContext()
	tokens, err := h.authService.VerifyEmailToken(ctx, token)
	if err != nil {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "Unathorized", err)
	}
	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Email Verified Successfull!",
		fiber.Map{"access_token": token},
	))
}

func (h *authHandlerImpl) Login(c *fiber.Ctx) error {
	var body LoginRequestDTO
	ctx := c.UserContext()
	if err := c.BodyParser(&body); err != nil {
		return apperr.BadRequest(apperr.CodeBadRequest, "Invalid Request", err)
	}
	tokens, err := h.authService.Login(ctx, &body)
	if err != nil {
		return err
	}
	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(
		httpx.NewHttpResponse(
			fiber.StatusOK,
			"Login successful",
			fiber.Map{"access_token": tokens.AccessToken},
		))
}

func (h *authHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
	oldRefreshToken := c.Cookies("refresh_token")
	if oldRefreshToken == "" {
		return  apperr.Unauthorized(apperr.CodeUnauthorized, "Unathorized", errors.New("Unathorized"))
	}

	tokens, err := h.authService.RefreshToken(ctx, oldRefreshToken)
	if err != nil {
		h.clearRefreshTokenCookie(c)
		return apperr.Unauthorized(apperr.CodeUnauthorized, "Unathorized", err)
	}

	h.setRefreshTokenCookie(c, tokens.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(httpx.NewHttpResponse(
		fiber.StatusOK,
		"Refresh Token successful",
		fiber.Map{"access_token": tokens.AccessToken},
	))
}

func (h *authHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()
	refreshToken := c.Cookies("refresh_token")
	if err := h.authService.Logout(ctx, refreshToken); err != nil {
		return err
	}
	h.clearRefreshTokenCookie(c)
	return c.Status(fiber.StatusBadRequest).JSON(httpx.NewHttpResponse[any](
		fiber.StatusBadRequest, 
		"logout Success", nil,
		))
}

func (h *authHandlerImpl) setRefreshTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
func (h *authHandlerImpl) clearRefreshTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
