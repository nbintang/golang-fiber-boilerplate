package auth

import (
	"context"
	"errors"
	"rest-fiber/config"
	"rest-fiber/internal/apperr"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/infra/cache"
	"rest-fiber/internal/infra/email"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/token"
	"rest-fiber/internal/user"
	"rest-fiber/pkg/password"
)

type authServiceImpl struct {
	userRepo     user.UserRepository
	tokenService token.Service
	emailService email.EmailService
	cacheService cache.Service
	env          config.Env
	logger       *infraapp.AppLogger
}

func NewAuthService(
	userRepo user.UserRepository,
	tokenService token.Service,
	emailService email.EmailService,
	cacheService cache.Service,
	env config.Env,
	logger *infraapp.AppLogger,
) AuthService {
	return &authServiceImpl{userRepo, tokenService, emailService, cacheService, env, logger}
}

func (s *authServiceImpl) Register(ctx context.Context, dto *RegisterRequestDTO) error {
	exists, err := s.userRepo.FindExistsByEmail(ctx, dto.Email)
	if err != nil {
		return apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	if exists {
		return apperr.Conflict(apperr.CodeConflict, "User Already Exists", errors.New("Email Already Exists"))
	}
	hashed, err := password.Hash(dto.Password)
	if err != nil {
		return apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	user := user.User{
		AvatarURL: dto.AvatarURL,
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  hashed,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	token, err := s.generateVerificationToken(user.ID.String())
	if err != nil {
		return apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	msg := s.env.FrontendURL + "/verify?token=" + token
	go s.sendEmail(user.Email, "Verification", msg)
	return nil
}

func (s *authServiceImpl) Login(ctx context.Context, dto *LoginRequestDTO) (TokensResponseDto, error) {
	user, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return TokensResponseDto{}, apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	if user == nil {
		return TokensResponseDto{}, apperr.NotFound(apperr.CodeNotFound, "User Not Found", errors.New("User Not Found"))
	}
	if err := password.Compare(dto.Password, user.Password); err != nil {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Password", err)
	}
	if user.IsEmailVerified == false {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Email Not Verified", errors.New("Email Not Verified"))
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}
func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (TokensResponseDto, error) {
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", err)
	}
	userID, _ := (*claims)["id"].(string)
	oldRT, ok := (*claims)["jti"].(string)
	if !ok || oldRT == "" || userID == "" {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", errors.New("Invalid Token"))
	}
	storedUserID, exists, err := s.cacheService.GetAndDel(ctx, keyRefresh+oldRT)
	if err != nil {
		return TokensResponseDto{}, apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	if !exists || storedUserID != userID {
		s.logger.Warnf("refresh token reuse detected for user %s with JTI %s", userID, oldRT)
		s.revokeAllUserTokens(ctx, userID)
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", errors.New("Invalid Token"))
	}
	s.blacklistAccessByRefreshJTI(ctx, oldRT)
	s.cacheService.Del(ctx, keyRTAccess+oldRT, keyRTAccessExp+oldRT)
	s.cacheService.SRem(ctx, keyUserTokens+userID, oldRT)
	user, err := s.userRepo.FindByIDWithRole(ctx, userID)
	if err != nil {
		return TokensResponseDto{}, apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	if user == nil {
		return TokensResponseDto{}, apperr.NotFound(apperr.CodeNotFound, "User Not Found", errors.New("User Not Found"))
	}
	if !user.IsEmailVerified {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Email Not Verified", errors.New("Email Not Verified"))
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}

func (s *authServiceImpl) VerifyEmailToken(ctx context.Context, verificationToken string) (TokensResponseDto, error) {
	claims, err := s.tokenService.VerifyToken(verificationToken, s.env.JWTVerificationSecret)
	if err != nil {
		return TokensResponseDto{}, apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", err)
	}
	userID := (*claims)["id"].(string)
	user, err := s.userRepo.FindByIDWithRole(ctx, userID)
	if err != nil {
		return TokensResponseDto{}, apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	if user == nil {
		return TokensResponseDto{}, apperr.NotFound(apperr.CodeNotFound, "User Not Found", errors.New("User Not Found"))
	}
	user.IsEmailVerified = true
	if err = s.userRepo.Update(ctx, user.ID.String(), user); err != nil {
		return TokensResponseDto{}, apperr.Internal(apperr.CodeInternal, "Internal Server Error", err)
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", errors.New("Invalid Token"))
	}
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", err)
	}
	userID, _ := (*claims)["id"].(string)
	rtJTI, _ := (*claims)["jti"].(string)
	if userID == "" || rtJTI == "" {
		return apperr.Unauthorized(apperr.CodeUnauthorized, "Invalid Token", errors.New("Invalid Token"))
	}
	s.blacklistAccessByRefreshJTI(ctx, rtJTI)
	s.cacheService.Del(ctx,
		keyRefresh+rtJTI,
		keyRTAccess+rtJTI,
		keyRTAccessExp+rtJTI,
	)
	s.cacheService.SRem(ctx, keyUserTokens+userID, rtJTI)
	return nil
}
