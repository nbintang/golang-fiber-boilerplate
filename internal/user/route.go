package user

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/http/middleware"
	"rest-fiber/internal/http/router"

	"github.com/gofiber/fiber/v2"
)

type UserRouteParams struct {
	router.InjectRouteParams
	UserHandler UserHandler
}
type userRouteImpl struct {
	userHandler UserHandler
}

func NewUserRoute(params UserRouteParams) router.ProtectedRoute {
	return &userRouteImpl{userHandler: params.UserHandler}
}
func (r *userRouteImpl) RegisterProtectedRoute(route fiber.Router) {
	users := route.Group("/users")
	users.Get("/me", r.userHandler.GetCurrentUserProfile)
	users.Patch("/me", r.userHandler.UpdateCurrentUser)
	users.Get("/", middleware.AllowRoleAccess(enums.Admin), r.userHandler.GetAllUsers)
	users.Get("/:id", middleware.AllowRoleAccess(enums.Admin), r.userHandler.GetUserByID)
}
