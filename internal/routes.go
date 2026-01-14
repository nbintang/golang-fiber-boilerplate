package internal

import (
	"rest-fiber/internal/http/router"
)

type RoutesIn struct {
	router.InjectRouteParams
	App       *Bootstrap
	Routes    []router.Route          `group:"public_routes"`
	Protected []router.ProtectedRoute `group:"protected_routes"`
}

func RegisterRoutes(in RoutesIn) {
	for _, r := range in.Routes {
		r.RegisterRoute(in.App.PublicRoute)
	}
	for _, r := range in.Protected {
		r.RegisterProtectedRoute(in.App.ProtectedRoute)
	}
}
