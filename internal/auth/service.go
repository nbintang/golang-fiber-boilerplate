package auth

import (
	"context"
	"errors"
	"rest-fiber/config"
	"rest-fiber/internal/enums"
	"rest-fiber/internal/infra/email"
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/rediscache"
	"rest-fiber/internal/infra/token"
	"rest-fiber/internal/user"
	"rest-fiber/pkg/password"
	"time"
)

type authServiceImpl struct {
	userRepo     user.UserRepository
	tokenService token.Service
	emailService email.EmailService
	redisService rediscache.Service
	env          config.Env
	logger       *infraapp.AppLogger
}

func NewAuthService(
	userRepo user.UserRepository,
	tokenService token.Service,
	emailService email.EmailService,
	redisService rediscache.Service,
	env config.Env,
	logger *infraapp.AppLogger,
) AuthService {
	return &authServiceImpl{userRepo, tokenService, emailService, redisService, env, logger}
}

func (s *authServiceImpl) Register(ctx context.Context, dto *RegisterRequestDTO) error {
	exists, err := s.userRepo.FindExistsByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("User Already Exist")
	}
	hashed, err := password.Hash(dto.Password)
	if err != nil {
		return err
	}
	user := user.User{
		AvatarURL: dto.AvatarURL,
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  hashed,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return err
	}
	token, err := s.generateVerificationToken(user.ID.String())
	if err != nil {
		return err
	}
	go s.sendEmail(user.Email, token, "Verification")
	return nil
}
func (s *authServiceImpl) sendEmail(emailAddr, token, subject string) {
	emailCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.emailService.SendEmail(emailCtx, email.Params{
		Subject: subject,
		Message: s.env.TargetURL + token,
		Reciever: email.Reciever{
			Email: emailAddr,
		}}); err != nil {
		s.logger.Error(err)
	}
}
func (s *authServiceImpl) Login(ctx context.Context, dto *LoginRequestDTO) (TokensResponseDto, error) {
	user, err := s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		return TokensResponseDto{}, err
	}
	if user == nil {
		return TokensResponseDto{}, errors.New("User Not Found")
	}
	if err := password.Compare(dto.Password, user.Password); err != nil {
		return TokensResponseDto{}, errors.New("Invalid Password")
	}
	if user.IsEmailVerified == false {
		return TokensResponseDto{}, errors.New("Email Not Verified")
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}
func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (TokensResponseDto, error) {
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return TokensResponseDto{}, err
	}
	userID, _ := (*claims)["id"].(string)
	oldRT, ok := (*claims)["jti"].(string)
	if !ok || oldRT == "" || userID == "" {
		return TokensResponseDto{}, errors.New("invalid token claims")
	}
	storedUserID, exists, err := s.redisService.GetAndDel(ctx, keyRefresh+oldRT)
	if err != nil {
		return TokensResponseDto{}, err
	}
	if !exists || storedUserID != userID {
		s.logger.Warnf("refresh token reuse detected for user %s with JTI %s", userID, oldRT)
		s.revokeAllUserTokens(ctx, userID)
		return TokensResponseDto{}, errors.New("token reuse detected")
	}
	s.blacklistAccessByRefreshJTI(ctx, oldRT)
	s.redisService.Del(ctx, keyRTAccess+oldRT, keyRTAccessExp+oldRT)
	s.redisService.SRem(ctx, keyUserTokens+userID, oldRT)
	user, err := s.userRepo.FindByIDWithRole(ctx, userID)
	if err != nil {
		return TokensResponseDto{}, err
	}
	if user == nil {
		return TokensResponseDto{}, errors.New("User Not Found")
	}
	if !user.IsEmailVerified {
		return TokensResponseDto{}, errors.New("Email Not Verified")
	}
	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}

func (s *authServiceImpl) VerifyEmailToken(ctx context.Context, verificationToken string) (TokensResponseDto, error) {
	claims, err := s.tokenService.VerifyToken(verificationToken, s.env.JWTVerificationSecret)
	if err != nil {
		return TokensResponseDto{}, err
	}
	userID := (*claims)["id"].(string)
	user, err := s.userRepo.FindByIDWithRole(ctx, userID)
	if err != nil {
		return TokensResponseDto{}, err
	}
	if user == nil {
		return TokensResponseDto{}, errors.New("User Not Found")
	}
	user.IsEmailVerified = true
	if err = s.userRepo.Update(ctx, user.ID.String(), user); err != nil {
		return TokensResponseDto{}, err
	}

	return s.generateTokens(ctx, user.ID.String(), user.Email, enums.EUserRoleType(user.Role))
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("no refresh token")
	}
	claims, err := s.tokenService.VerifyToken(refreshToken, s.env.JWTRefreshSecret)
	if err != nil {
		return err
	}
	userID, _ := (*claims)["id"].(string)
	rtJTI, _ := (*claims)["jti"].(string)
	if userID == "" || rtJTI == "" {
		return errors.New("invalid token claims")
	}
	s.blacklistAccessByRefreshJTI(ctx, rtJTI)
	s.redisService.Del(ctx,
		keyRefresh+rtJTI,
		keyRTAccess+rtJTI,
		keyRTAccessExp+rtJTI,
	)
	s.redisService.SRem(ctx, keyUserTokens+userID, rtJTI)
	return nil
}
