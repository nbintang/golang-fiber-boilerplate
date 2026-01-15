package user

import (
	"context"
	"fmt"
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


func (s *userServiceImpl) getListVersion(ctx context.Context, versionKey string) int64 {
	raw, err := s.cacheService.Get(ctx, versionKey)
	if err == nil && raw != "" {
		var v int64
		_, scanErr := fmt.Sscan(raw, &v)
		if scanErr == nil && v > 0 {
			return v
		}
	}
	_ = s.cacheService.Set(ctx, versionKey, "1", 0)
	return 1
}


