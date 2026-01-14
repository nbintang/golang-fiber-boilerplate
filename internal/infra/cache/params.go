package cache

import (
	"time"

	"github.com/gofiber/storage"
)

type ThrottleParams struct {
	MaxLimit int
	Storage  storage.Storage
	Expiration time.Duration
	Prefix string
}