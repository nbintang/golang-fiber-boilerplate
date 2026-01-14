package user

import (
	"context"
	"fmt"
)

func (s *userServiceImpl) getListVersion(ctx context.Context, versionKey string) int64 {
	raw, err := s.cacheService.Get(ctx, versionKey)
	if err == nil && raw != "" {
		var v int64
		_, scanErr := fmt.Sscan(raw, &v)
		if scanErr == nil && v > 0 {
			return v
		}
	}
	_ = s.cacheService.Set(ctx, versionKey, "1", 0) // ttl 0 = no-expire (umum)
	return 1
}
