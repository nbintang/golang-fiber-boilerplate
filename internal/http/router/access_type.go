package router

type AccessType string

const (
	RouteProtected AccessType = `group:"protected_routes"`
	RoutePublic    AccessType = `group:"public_routes"`
)
