package enums

type EAccessType string

const (
	RouteProtected EAccessType = `group:"protected_routes"`
	RoutePublic    EAccessType = `group:"public_routes"`
)
