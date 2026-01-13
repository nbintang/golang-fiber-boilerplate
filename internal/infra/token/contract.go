package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(params *GenerateTokenParams, envStr string, duration time.Duration) (string, error)
	VerifyToken(tokenStr string, envStr string) (*jwt.MapClaims, error)
	RemainingTTLFromAccessToken(accessToken string) (time.Duration, error)
}
