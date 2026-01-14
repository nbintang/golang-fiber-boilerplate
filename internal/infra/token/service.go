package token

import (
	"errors"
	"rest-fiber/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type serviceImpl struct{}

func NewService(env config.Env) Service {
	return &serviceImpl{}
}

func (s *serviceImpl) GenerateToken(params *GenerateTokenParams, envStr string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := jwt.MapClaims{
		"id":   params.ID,
		"exp":  expirationTime.Unix(),
		"type": string(params.Type),
		"iat":  time.Now().Unix(),
		"role": params.Role,
	}
	if params.JTI != "" {
		claims["jti"] = params.JTI
	}
	if params.Email != "" {
		claims["email"] = params.Email
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(envStr))
}

func (s *serviceImpl) VerifyToken(tokenStr string, envStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(envStr), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return &claims, nil
}

func (s *serviceImpl) RemainingTTLFromAccessToken(accessToken string) (time.Duration, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return 0, errors.New("invalid exp claim")
	}

	expTime := time.Unix(int64(expFloat), 0)
	ttl := time.Until(expTime)

	if ttl < 0 {
		ttl = 0
	}

	return ttl, nil
}
