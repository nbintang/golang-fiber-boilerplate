package auth

import (
	"rest-fiber/internal/http/router"
	"rest-fiber/internal/infra/rediscache"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthRouteParams struct {
	router.RouteParams
	AuthHandler  AuthHandler
	RedisService rediscache.Service
}

type authRouteImpl struct {
	authHandler  AuthHandler
	redisService rediscache.Service
}

func NewAuthRoute(params AuthRouteParams) router.Route {
	return &authRouteImpl{authHandler: params.AuthHandler, redisService: params.RedisService}
}
func (r *authRouteImpl) RegisterRoute(api fiber.Router) {
	auth := api.Group("/auth")

	redisStarage := r.redisService.GetStorage()
	storageParams := rediscache.ThrottleParams{MaxLimit: 5, Storage: redisStarage, Expiration: 1 * time.Minute}
	auth.Post("/register", rediscache.Throttle(storageParams), r.authHandler.Register)
	auth.Post("/verify", rediscache.Throttle(storageParams), r.authHandler.VerifyEmail)
	auth.Post("/login", rediscache.Throttle(storageParams), r.authHandler.Login)
	auth.Delete("/logout", rediscache.Throttle(storageParams), r.authHandler.Logout)
	auth.Post("/refresh-token", rediscache.Throttle(storageParams), r.authHandler.RefreshToken)
}
