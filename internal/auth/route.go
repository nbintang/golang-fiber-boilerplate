	package auth

	import (
		"rest-fiber/internal/http/router"
		"rest-fiber/internal/infra/cache"
		"time"

		"github.com/gofiber/fiber/v2"
	)

	type AuthRouteParams struct {
		router.InjectRouteParams
		AuthHandler  AuthHandler
		CacheService cache.Service
	}

	type authRouteImpl struct {
		authHandler  AuthHandler
		cacheService cache.Service
	}

	func NewAuthRoute(params AuthRouteParams) router.Route {
		return &authRouteImpl{authHandler: params.AuthHandler, cacheService: params.CacheService}
	}
	func (r *authRouteImpl) RegisterRoute(api fiber.Router) {
		auth := api.Group("/auth")

		redisStarage := r.cacheService.GetStorage()
		storageParams := cache.ThrottleParams{MaxLimit: 5, Storage: redisStarage, Expiration: 1 * time.Minute}
		auth.Post("/register", cache.Throttle(storageParams), r.authHandler.Register)
		auth.Post("/verify", cache.Throttle(storageParams), r.authHandler.VerifyEmail)
		auth.Post("/login", cache.Throttle(storageParams), r.authHandler.Login)
		auth.Delete("/logout", cache.Throttle(storageParams), r.authHandler.Logout)
		auth.Post("/refresh-token", cache.Throttle(storageParams), r.authHandler.RefreshToken)
	}
