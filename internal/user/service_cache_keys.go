package user

import (
	"fmt"
	"time"
)

const (
	userByIDTTL     time.Duration = 10 * time.Minute
	usersListTTL    time.Duration = 60 * time.Second
	usersVersionKey string        = "users:version"
)

type usersListCache struct {
	Total int64             `json:"total"`
	Items []UserResponseDTO `json:"items"`
}

func userByIDKey(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func usersListKey(version int64, limit, offset int) string {
	return fmt.Sprintf("users:v%d:limit:%d:offset:%d", version, limit, offset)
}
