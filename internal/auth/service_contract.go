package auth

import "context"

type AuthService interface {
	Register(ctx context.Context, dto *RegisterRequestDTO) error
	VerifyEmailToken(ctx context.Context, token string) (TokensResponseDto, error)
	Login(ctx context.Context, dto *LoginRequestDTO) (TokensResponseDto, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokensResponseDto, error)
	Logout(ctx context.Context, refreshToken string) error
}
