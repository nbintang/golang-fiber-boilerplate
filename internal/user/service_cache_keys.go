package user

import "time"

const (
	userByIDTTL     time.Duration = 10 * time.Minute
	usersListTTL    time.Duration = 60 * time.Second
	usersVersionKey string        = "users:version"
)

